import { Component, createEffect, onCleanup, For, JSX } from "solid-js";
import UserInput from "./UserInput";
import UserMessage from "./UserMessage";
import BotMessage from "./BotMessage";
import { TerminalMessageProps } from "types";

export interface TerminalDialogProps {
  isOpen: boolean;
  onClose: () => void;
  onSubmit: (query: string) => Promise<void>;
  chatHistory: TerminalMessageProps[];
  isLoading: boolean;
  children?: JSX.Element;
}

const TerminalDialog: Component<TerminalDialogProps> = (props) => {
  let dialogRef: HTMLDialogElement | undefined;

  createEffect(() => {
    if (dialogRef) {
      if (props.isOpen && !dialogRef.open) {
        dialogRef.showModal();
      } else if (!props.isOpen && dialogRef.open) {
        dialogRef.close();
      }
    }
  });

  const handleKeyDown = (e: KeyboardEvent) => {
    if (e.key === "Escape" && props.isOpen) {
      props.onClose();
    }
  };

  const handleBackdropClick = (e: MouseEvent) => {
    if (dialogRef && e.target === dialogRef) {
      props.onClose();
    }
  };

  // Add event listeners
  createEffect(() => {
    if (props.isOpen) {
      document.addEventListener("keydown", handleKeyDown);
      if (dialogRef) {
        dialogRef.addEventListener("click", handleBackdropClick);
      }
    }

    onCleanup(() => {
      document.removeEventListener("keydown", handleKeyDown);
      if (dialogRef) {
        dialogRef.removeEventListener("click", handleBackdropClick);
      }
    });
  });

  return (
    <dialog
      ref={dialogRef}
      class="!fixed top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 m-0 w-4/5 aspect-square max-w-3xl bg-black/85 text-blue-300 p-6 rounded-2xl border-4 border-blue-300 focus:outline-none"
      style={{
        "font-family": "'VT323', 'Courier New', monospace",
      }}
    >
      <div class="flex flex-col h-full">
        <div class="flex justify-end items-center mb-4">
          <button
            onClick={props.onClose}
            class="text-blue-300 hover:text-blue-200 focus:outline-none"
            aria-label="Close dialog"
          >
            X
          </button>
        </div>
        <div class="flex-grow overflow-auto terminal-content mb-4">
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
            autoFocus={props.isOpen}
          />
        </div>
      </div>
    </dialog>
  );
};

export default TerminalDialog;
