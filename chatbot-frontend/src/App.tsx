import type { Component } from "solid-js";
import { createSignal, JSX, onMount, onCleanup } from "solid-js";
import TerminalButton from "./components/TerminalButton";
import TerminalDialog from "./components/TerminalDialog";
import { createChatBuffer } from "./components/ChatBuffer";
import { Buffer } from "./types";

const App: Component = () => {
  const [isDialogOpen, setIsDialogOpen] = createSignal(false);
  const [buffers, setBuffers] = createSignal<Buffer[]>([]);
  const [activeBufferId, setActiveBufferId] = createSignal<string>("");
  const [navbar, setNavbar] = createSignal<JSX.Element | null>(null);

  const createChatBufferInstance = () => {
    const chatBufferId = "chat";

    const chatBuffer = createChatBuffer({
      id: chatBufferId,
      name: "Chat",
      onExit: closeDialog,
      setNavbar,
    });

    setBuffers([chatBuffer]);
    setActiveBufferId(chatBufferId);

    return chatBuffer;
  };

  const getActiveBuffer = () => {
    let currentBuffer = buffers().find((b) => b.id === activeBufferId());

    if (!currentBuffer) {
      currentBuffer = createChatBufferInstance();
    }

    return currentBuffer;
  };

  const openDialog = () => setIsDialogOpen(true);
  const closeDialog = () => setIsDialogOpen(false);

  const handleKeyDown = (event: KeyboardEvent) => {
    if (
      event.key === "/" &&
      !isDialogOpen() &&
      !(event.target instanceof HTMLInputElement) &&
      !(event.target instanceof HTMLTextAreaElement)
    ) {
      event.preventDefault();
      openDialog();
    }
  };

  onMount(() => {
    document.addEventListener("keydown", handleKeyDown);

    onCleanup(() => {
      document.removeEventListener("keydown", handleKeyDown);
    });
  });

  return (
    <>
      <TerminalButton onClick={openDialog} />
      {isDialogOpen() && (
        <TerminalDialog
          isOpen={isDialogOpen()}
          onClose={closeDialog}
          activeBuffer={getActiveBuffer()}
          navbar={navbar()}
        />
      )}
    </>
  );
};

export default App;
