import { Component, For, createSignal, createEffect, JSX } from "solid-js";
import {
  Buffer,
  ChatBufferProps,
  TerminalMessageProps,
  Source,
} from "../types";
import UserMessage from "./UserMessage";
import BotMessage from "./BotMessage";
import UserInput from "./UserInput";
import Navbar from "./Navbar";
import { filePathToURL } from "../utils/siteUtils";

const ChatBuffer: Component<{
  id: string;
  chatHistory?: TerminalMessageProps[];
  isLoading?: boolean;
  apiUrl?: string;
  onExit?: () => void;
  setNavbar?: (navbar: JSX.Element) => void;
}> = (props) => {
  const defaultApiUrl =
    import.meta.env.VITE_API_URL || "http://localhost:8181/chat";

  const storageKey = `chat_history_${props.id}`;
  const sourcesStorageKey = `chat_sources_${props.id}`;
  const sourceIndexStorageKey = `chat_source_index_${props.id}`;

  const loadChatHistory = (): TerminalMessageProps[] => {
    if (typeof window === "undefined") return props.chatHistory || [];

    try {
      const savedHistory = localStorage.getItem(storageKey);
      if (savedHistory) {
        return JSON.parse(savedHistory);
      }
    } catch (error) {
      console.error("Error loading chat history from localStorage:", error);
    }

    return (
      props.chatHistory || [
        {
          type: "bot",
          content:
            "Hello Professor\nWould you like to learn about Epistemic Technology?",
        },
      ]
    );
  };

  const loadSources = (): Source[] => {
    if (typeof window === "undefined") return [];

    try {
      const savedSources = localStorage.getItem(sourcesStorageKey);
      if (savedSources) {
        return JSON.parse(savedSources);
      }
    } catch (error) {
      console.error("Error loading sources from localStorage:", error);
    }

    return [];
  };

  const loadSourceIndex = (): number => {
    if (typeof window === "undefined") return 0;

    try {
      const savedSourceIndex = localStorage.getItem(sourceIndexStorageKey);
      if (savedSourceIndex) {
        return JSON.parse(savedSourceIndex);
      }
    } catch (error) {
      console.error("Error loading source index from localStorage:", error);
    }

    return 0;
  };

  const [chatHistory, setChatHistory] = createSignal<TerminalMessageProps[]>(
    loadChatHistory()
  );
  const [isLoading, setIsLoading] = createSignal(props.isLoading || false);
  const [apiUrl] = createSignal(props.apiUrl || defaultApiUrl);
  const [sources, setSources] = createSignal<Source[]>(loadSources());
  const [sourceIndex, setSourceIndex] = createSignal(loadSourceIndex());

  const addSource = (source: Source): Source => {
    if (sources().length < 9) {
      source.index = sources().length + 1;
      setSources((prev) => [...prev, source]);
      setSourceIndex((prev) => prev + 1);
      console.log(sourceIndex());
      return source;
    } else {
      setSources((prev) => {
        source.index = (sourceIndex() + 1) % 9;
        const newSources = [...prev];
        newSources[sourceIndex()] = source;
        setSourceIndex((prev) => (prev + 1) % 9);
        console.log(sourceIndex());
        return newSources;
      });
      return source;
    }
  };

  createEffect(() => {
    const currentHistory = chatHistory();
    if (typeof window !== "undefined") {
      try {
        localStorage.setItem(storageKey, JSON.stringify(currentHistory));
      } catch (error) {
        console.error("Error saving chat history to localStorage:", error);
      }
    }
  });

  createEffect(() => {
    const currentSources = sources();
    const currentSourceIndex = sourceIndex();

    if (typeof window !== "undefined") {
      try {
        localStorage.setItem(sourcesStorageKey, JSON.stringify(currentSources));
        localStorage.setItem(
          sourceIndexStorageKey,
          JSON.stringify(currentSourceIndex)
        );
      } catch (error) {
        console.error("Error saving sources to localStorage:", error);
      }
    }
  });

  const helpCommand = () => {
    const helpText = `Available commands:
    :help - Show this help
    :about - About the chatbot
    :contact - Contact information
    :clear - Clear the chat history
    :reset - Reset the chat to initial state
    :exit - Exit the chat
    :sources - Show list of sources
    :1-9 - Open a source in the browser
    `;
    setChatHistory((prev) => [
      ...prev,
      {
        type: "bot",
        content: helpText,
      },
    ]);
  };

  const aboutCommand = () => {
    const aboutText = `Joshua is a representative of Epistemic Technology.

    He can answer your questions about us by searching our site for relevant information and summarizing what he finds.

    We have asked Joshua to be helpful and professional, and to only answer questions related to what he finds on the site. If you manage to get him off topic, or to say something unbecoming of a representative of Epistemic Technology, please let us know!

    Joshua also has some other abilities, which you can see using the :help command.`;

    const aboutSource: Source = {
      Title: "About our chatbot",
      FilePath:
        "/Users/mikethicke/github/epistemic.technology/site/content/about-our-chatbot.md",
      Author: "",
      ID: -1,
      Content: "",
      PublicationDate: "",
    };

    setChatHistory((prev) => [
      ...prev,
      {
        type: "bot",
        content: aboutText,
        sources: [addSource(aboutSource)],
      },
    ]);
  };

  const contactCommand = () => {
    const contactFilePath =
      "/Users/mikethicke/github/epistemic.technology/site/content/contact.md";
    const contactURL = filePathToURL(contactFilePath);
    window.location.href = contactURL;
  };

  const sourceCommand = () => {
    for (const source of sources()) {
      setChatHistory((prev) => [
        ...prev,
        {
          type: "bot",
          content: `${source.index}: ${source.Title} - ${filePathToURL(
            source.FilePath
          )}`,
        },
      ]);
    }
  };

  const goSourceCommand = (index: number) => {
    const sourcesArray = sources();
    if (
      index >= 0 &&
      index < sourcesArray.length &&
      sourcesArray[index].FilePath
    ) {
      window.location.href = filePathToURL(sourcesArray[index].FilePath);
    } else {
      setChatHistory((prev) => [
        ...prev,
        {
          type: "bot",
          content: `No source available at index ${index + 1}.`,
        },
      ]);
    }
  };

  const handleChatSubmit = async (query: string) => {
    setChatHistory((prev) => [
      ...prev,
      {
        type: "user",
        content: query,
      },
    ]);

    // Scroll to the bottom after user message is added
    setTimeout(() => {
      const userInputContainer = document.querySelector(
        ".user-input-container"
      );
      if (userInputContainer) {
        userInputContainer.scrollIntoView({
          behavior: "smooth",
          block: "end",
          inline: "nearest",
        });
      }
    }, 100);

    if (query.startsWith(":")) {
      switch (query.slice(1)) {
        case "help":
          helpCommand();
          return;
        case "about":
          aboutCommand();
          return;
        case "contact":
          contactCommand();
          return;
        case "clear":
          setChatHistory([]);
          return;
        case "reset":
          localStorage.removeItem(storageKey);
          localStorage.removeItem(sourcesStorageKey);
          localStorage.removeItem(sourceIndexStorageKey);
          setChatHistory(loadChatHistory());
          setSources([]);
          setSourceIndex(0);
          return;
        case "exit":
          if (props.onExit) {
            props.onExit();
          }
          return;
        case "sources":
          sourceCommand();
          return;
        case "1":
        case "2":
        case "3":
        case "4":
        case "5":
        case "6":
        case "7":
        case "8":
        case "9":
          const index = parseInt(query.slice(1)) - 1;
          goSourceCommand(index);
          return;
      }
    }

    try {
      setIsLoading(true);

      const response = await fetch(apiUrl(), {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ query }),
      });

      if (!response.ok) {
        throw new Error(`API error: ${response.status}`);
      }

      const data = await response.json();
      const sources = [];
      for (const source of data.sources) {
        sources.push(addSource(source));
      }

      setChatHistory((prev) => [
        ...prev,
        {
          type: "bot",
          content: data.response,
          sources: sources,
        },
      ]);
    } catch (error) {
      console.error("Error submitting chat:", error);

      setChatHistory((prev) => [
        ...prev,
        {
          type: "bot",
          content: "Sorry, an error occurred while processing your request.",
        },
      ]);
    } finally {
      setIsLoading(false);

      // Scroll to the bottom after bot message is added
      setTimeout(() => {
        const userInputContainer = document.querySelector(
          ".user-input-container"
        );
        if (userInputContainer) {
          userInputContainer.scrollIntoView({
            behavior: "smooth",
            block: "end",
            inline: "nearest",
          });
        }
      }, 100);
    }
  };

  const handleCommandExecute = (command: string) => {
    handleChatSubmit(command);
  };

  if (props.setNavbar) {
    const chatNavbar = (
      <Navbar
        commands={[
          { name: ":help", command: ":help" },
          { name: ":about", command: ":about" },
          { name: ":contact", command: ":contact" },
          { name: ":exit", command: ":exit" },
        ]}
        onCommandExecute={handleCommandExecute}
      />
    );
    props.setNavbar(chatNavbar);
  }

  return (
    <div class="flex-grow terminal-content mb-4 px-6">
      <For each={chatHistory()}>
        {(message) => (
          <div>
            {message.type === "user" && (
              <UserMessage content={message.content} type="user" />
            )}
            {message.type === "bot" && (
              <BotMessage
                content={message.content}
                sources={message.sources}
                type="bot"
              />
            )}
          </div>
        )}
      </For>
      <UserInput
        onSubmit={handleChatSubmit}
        isLoading={isLoading()}
        autoFocus={true}
      />
    </div>
  );
};

// Factory function that creates a Buffer containing a chat interface
export const createChatBuffer = (
  props: ChatBufferProps & {
    id: string;
    name: string;
    setNavbar?: (navbar: JSX.Element) => void;
  }
): Buffer => {
  return {
    id: props.id,
    name: props.name,
    render: () => (
      <ChatBuffer
        id={props.id}
        chatHistory={props.chatHistory}
        isLoading={props.isLoading}
        apiUrl={props.apiUrl}
        onExit={props.onExit}
        setNavbar={props.setNavbar}
      />
    ),
  };
};
