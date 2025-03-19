import type { Accessor, Component, Setter } from "solid-js";
import { createSignal } from "solid-js";
import TerminalButton from "./components/TerminalButton";
import TerminalDialog from "./components/TerminalDialog";
import {
  BotMessageProps,
  TerminalMessageProps,
  UserMessageProps,
} from "./types";

const App: Component = () => {
  const [isDialogOpen, setIsDialogOpen] = createSignal(false);
  const [apiUrl, setApiUrl] = createSignal(
    import.meta.env.VITE_API_URL || "http://localhost:8181/chat"
  );
  const [chatHistory, setChatHistory] = createSignal<TerminalMessageProps[]>([
    {
      type: "bot",
      content:
        "Hello Professor\nWould you like to learn about Epistemic Technology?",
    },
  ]);
  const [isLoading, setIsLoading] = createSignal(false);

  const openDialog = () => setIsDialogOpen(true);
  const closeDialog = () => setIsDialogOpen(false);

  const handleChatSubmit = async (query: string) => {
    setIsLoading(true);

    const newMessage: UserMessageProps = {
      type: "user",
      content: query,
    };
    setChatHistory([...chatHistory(), newMessage]);

    if (query.startsWith(":")) {
      switch (query.slice(1)) {
        case "help":
          helpCommand(chatHistory, setChatHistory);
          break;
        case "clear":
          setChatHistory([]);
          break;
        case "exit":
          setIsDialogOpen(false);
          break;
      }
      setIsLoading(false);
      return;
    }

    try {
      const historyString = JSON.stringify(
        chatHistory().map((msg) => ({
          content: msg.content,
        }))
      );

      const response = await fetch(apiUrl(), {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({
          query: query,
          history: historyString,
        }),
      });

      if (!response.ok) {
        throw new Error(`API error: ${response.status}`);
      }

      const data = await response.json();

      const responseMessage: BotMessageProps = {
        type: "bot",
        content: data.response,
        sources: data.sources,
      };

      setChatHistory([...chatHistory(), responseMessage]);
    } catch (error) {
      console.error("Error sending message:", error);
      const updatedHistory = [...chatHistory()];
      updatedHistory[updatedHistory.length - 1] = {
        ...updatedHistory[updatedHistory.length - 1],
        content: "Error: Could not connect to the server.",
      };
      setChatHistory(updatedHistory);
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <>
      <TerminalButton onClick={openDialog} />
      <TerminalDialog
        isOpen={isDialogOpen()}
        onClose={closeDialog}
        onSubmit={handleChatSubmit}
        chatHistory={chatHistory()}
        isLoading={isLoading()}
      />
    </>
  );
};

function helpCommand(
  chatHistory: Accessor<TerminalMessageProps[]>,
  setChatHistory: Setter<TerminalMessageProps[]>
) {
  const helpText = `Available commands:
  :help - Show this help
  :clear - Clear the chat history
  :exit - Exit the chat
  `;
  setChatHistory([
    ...chatHistory(),
    {
      type: "bot",
      content: helpText,
    },
  ]);
}

export default App;
