import {
  Component,
  createEffect,
  onCleanup,
  JSX,
  createSignal,
} from "solid-js";
import { Buffer } from "../types";

export interface TerminalDialogProps {
  isOpen: boolean;
  onClose: () => void;
  activeBuffer: Buffer;
  executeCommand: (command: string) => void;
  children?: JSX.Element;
}

const TerminalDialog: Component<TerminalDialogProps> = (props) => {
  let dialogRef: HTMLDialogElement | undefined;
  let terminalContentRef: HTMLDivElement | undefined;
  const [inputElement, setInputElement] = createSignal<HTMLInputElement | null>(
    null
  );

  // Focus the input element
  const focusInput = () => {
    // Find the input element within the active buffer
    if (!inputElement()) {
      if (terminalContentRef) {
        const input = terminalContentRef.querySelector(
          "input"
        ) as HTMLInputElement;
        if (input) {
          setInputElement(input);
        }
      }
    }

    if (inputElement()) {
      inputElement()!.focus();
    }
  };

  createEffect(() => {
    if (dialogRef) {
      if (props.isOpen && !dialogRef.open) {
        dialogRef.showModal();
        // Focus the input after the dialog opens
        setTimeout(focusInput, 50);
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

  // Handle clicks within the dialog to ensure input stays focused
  const handleDialogClick = (e: MouseEvent) => {
    // Only focus if we're not clicking on a button or another interactive element
    const target = e.target as HTMLElement;
    if (
      target.tagName !== "BUTTON" &&
      !target.closest("button") &&
      !target.closest("a") &&
      target !== inputElement()
    ) {
      e.preventDefault();
      focusInput();
    }
  };

  // Add event listeners
  createEffect(() => {
    if (props.isOpen) {
      document.addEventListener("keydown", handleKeyDown);
      if (dialogRef) {
        dialogRef.addEventListener("click", handleBackdropClick);
        dialogRef.addEventListener("click", handleDialogClick);
      }

      // Find and set the input element
      setTimeout(focusInput, 50);
    }

    onCleanup(() => {
      document.removeEventListener("keydown", handleKeyDown);
      if (dialogRef) {
        dialogRef.removeEventListener("click", handleBackdropClick);
        dialogRef.removeEventListener("click", handleDialogClick);
      }
    });
  });

  // Refocus input when activeBuffer changes
  createEffect(() => {
    if (props.isOpen && props.activeBuffer) {
      // Reset the input element reference when buffer changes
      setInputElement(null);
      setTimeout(focusInput, 50);
    }
  });

  return (
    <dialog
      ref={dialogRef}
      class="!fixed top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 m-0 w-4/5 aspect-square max-w-3xl bg-black/85 text-blue-300 p-0 rounded-2xl border-4 border-blue-300 focus:outline-none overflow-hidden"
      style={{
        "font-family": "'VT323', 'Courier New', monospace",
      }}
    >
      <div class="flex flex-col h-full">
        <div class="flex justify-end items-center mb-4 px-6 pt-6">
          <button
            onClick={props.onClose}
            class="text-blue-300 hover:text-blue-200 focus:outline-none"
            aria-label="Close dialog"
          >
            X
          </button>
        </div>

        {/* Render the active buffer */}
        <div
          ref={terminalContentRef}
          class="flex-1 overflow-auto terminal-dialog-content px-0 pr-4 mx-2"
        >
          {props.activeBuffer.render()}
        </div>

        {/* Navigation Bar */}
        <div class="bg-blue-300 text-black p-2 flex justify-center space-x-6 w-full mt-auto">
          <button
            class="font-bold hover:underline focus:outline-none"
            onClick={() => props.executeCommand(":help")}
          >
            :help
          </button>
          <button
            class="font-bold hover:underline focus:outline-none"
            onClick={() => props.executeCommand(":about")}
          >
            :about
          </button>
          <button
            class="font-bold hover:underline focus:outline-none"
            onClick={() => props.executeCommand(":contact")}
          >
            :contact
          </button>
          <button
            class="font-bold hover:underline focus:outline-none"
            onClick={() => props.executeCommand(":exit")}
          >
            :exit
          </button>
        </div>
      </div>
    </dialog>
  );
};

export default TerminalDialog;
