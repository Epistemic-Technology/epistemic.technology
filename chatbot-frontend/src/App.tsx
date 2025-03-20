import type { Component } from "solid-js";
import { createSignal } from "solid-js";
import TerminalButton from "./components/TerminalButton";
import TerminalDialog from "./components/TerminalDialog";
import { createChatBuffer } from "./components/ChatBuffer";
import { Buffer } from "./types";

const App: Component = () => {
  const [isDialogOpen, setIsDialogOpen] = createSignal(false);
  const [buffers, setBuffers] = createSignal<Buffer[]>([]);
  const [activeBufferId, setActiveBufferId] = createSignal<string>("");

  // Create the chat buffer
  const createChatBufferInstance = () => {
    const chatBufferId = "chat";

    // Create the buffer using the factory function
    const chatBuffer = createChatBuffer({
      id: chatBufferId,
      name: "Chat",
    });

    setBuffers([chatBuffer]);
    setActiveBufferId(chatBufferId);

    return chatBuffer;
  };

  // Get the active buffer
  const getActiveBuffer = () => {
    let currentBuffer = buffers().find((b) => b.id === activeBufferId());

    // If no active buffer exists, create the chat buffer
    if (!currentBuffer) {
      currentBuffer = createChatBufferInstance();
    }

    return currentBuffer;
  };

  const openDialog = () => setIsDialogOpen(true);
  const closeDialog = () => setIsDialogOpen(false);

  // Execute a command in the active buffer
  const executeCommand = async (command: string) => {
    const activeBuffer = getActiveBuffer();
    if (activeBuffer && activeBuffer.handleCommand) {
      await activeBuffer.handleCommand(command);

      // Special case for exit command
      if (command === ":exit") {
        closeDialog();
      }
    }
  };

  return (
    <>
      <TerminalButton onClick={openDialog} />
      <TerminalDialog
        isOpen={isDialogOpen()}
        onClose={closeDialog}
        activeBuffer={getActiveBuffer()}
        executeCommand={executeCommand}
      />
    </>
  );
};

export default App;
