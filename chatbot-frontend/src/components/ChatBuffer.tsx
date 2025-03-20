import { Component, For, createSignal } from "solid-js";
import { Buffer, ChatBufferProps, TerminalMessageProps } from "../types";
import UserMessage from "./UserMessage";
import BotMessage from "./BotMessage";
import UserInput from "./UserInput";

// The component that renders the chat UI
const ChatBufferUI: Component<{
  chatHistory: TerminalMessageProps[];
  isLoading: boolean;
  onSubmit: (query: string) => Promise<void>;
}> = (props) => {
  return (
    <div class="flex-grow terminal-content mb-4 px-6">
      <For each={props.chatHistory}>
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
        onSubmit={props.onSubmit}
        isLoading={props.isLoading}
        autoFocus={true}
      />
    </div>
  );
};

// Factory function that creates a Buffer containing a chat interface
export const createChatBuffer = (
  props: ChatBufferProps & { id: string; name: string }
): Buffer => {
  const defaultApiUrl =
    import.meta.env.VITE_API_URL || "http://localhost:8181/chat";

  // Local state for chat functionality
  const [chatHistory, setChatHistory] = createSignal<TerminalMessageProps[]>(
    props.chatHistory || [
      {
        type: "bot",
        content:
          "Hello Professor\nWould you like to learn about Epistemic Technology?",
      },
    ]
  );
  const [isLoading, setIsLoading] = createSignal(props.isLoading || false);
  const [apiUrl, setApiUrl] = createSignal(props.apiUrl || defaultApiUrl);

  // Chat submission handler
  const handleChatSubmit = async (query: string) => {
    // Handle special commands
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
        case "exit":
          return;
      }
    }

    try {
      setIsLoading(true);

      // Add user message to chat history
      setChatHistory((prev) => [
        ...prev,
        {
          type: "user",
          content: query,
        },
      ]);

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

      // Add bot response to chat history
      setChatHistory((prev) => [
        ...prev,
        {
          type: "bot",
          content: data.response,
          sources: data.sources,
        },
      ]);
    } catch (error) {
      console.error("Error submitting chat:", error);

      // Add error message to chat history
      setChatHistory((prev) => [
        ...prev,
        {
          type: "bot",
          content: "Sorry, an error occurred while processing your request.",
        },
      ]);
    } finally {
      setIsLoading(false);
    }
  };

  // Command handlers
  const helpCommand = () => {
    const helpText = `
    Available commands:
    :help - Show this help
    :about - About Epistemic Technology
    :contact - Contact information
    :clear - Clear the chat history
    :exit - Exit the chat
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
    const aboutText = `About Epistemic Technology

Epistemic Technology is a field that focuses on the development of tools and systems that help us understand, create, and share knowledge more effectively.

Our mission is to build technologies that enhance human cognition and improve our ability to make sense of complex information landscapes.
`;
    setChatHistory((prev) => [
      ...prev,
      {
        type: "bot",
        content: aboutText,
      },
    ]);
  };

  const contactCommand = () => {
    const contactText = `Contact Information

Email: info@epistemic.technology
Website: https://epistemic.technology
Twitter: @EpistemicTech

Feel free to reach out with questions, ideas, or collaboration proposals.
`;
    setChatHistory((prev) => [
      ...prev,
      {
        type: "bot",
        content: contactText,
      },
    ]);
  };

  // If onSubmit prop is provided, use it, otherwise use local handler
  const submitHandler = props.onSubmit || handleChatSubmit;

  // Create and return a Buffer
  return {
    id: props.id,
    name: props.name,
    render: () => (
      <ChatBufferUI
        chatHistory={chatHistory()}
        isLoading={isLoading()}
        onSubmit={submitHandler}
      />
    ),
    handleCommand: submitHandler,
  };
};

// For backwards compatibility and simpler component usage
const ChatBuffer: Component<ChatBufferProps & { id: string; name: string }> = (
  props
) => {
  const buffer = createChatBuffer(props);
  return buffer.render();
};

export default ChatBuffer;
