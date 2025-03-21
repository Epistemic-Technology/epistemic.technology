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
  navbar?: JSX.Element;
  setNavbar?: (navbar: JSX.Element) => void;
}

const TerminalDialog: Component<TerminalDialogProps> = (props) => {
  let dialogRef: HTMLDialogElement | undefined;
  let terminalContentRef: HTMLDivElement | undefined;
  const [inputElement, setInputElement] = createSignal<HTMLInputElement | null>(
    null
  );

  const focusInput = () => {
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

  const handleDialogClick = (e: MouseEvent) => {
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

  createEffect(() => {
    if (props.isOpen) {
      document.addEventListener("keydown", handleKeyDown);
      if (dialogRef) {
        dialogRef.addEventListener("click", handleBackdropClick);
        dialogRef.addEventListener("click", handleDialogClick);
      }

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

  createEffect(() => {
    if (props.isOpen && props.activeBuffer) {
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
        <div class="flex justify-end items-center mb-1 px-6 pt-2">
          <button
            onClick={props.onClose}
            class="text-blue-300 hover:text-blue-200 focus:outline-none"
            aria-label="Close dialog"
          >
            X
          </button>
        </div>

        <div
          ref={terminalContentRef}
          class="flex-1 overflow-auto terminal-dialog-content px-0 pr-4 mx-2"
        >
          {props.activeBuffer.render()}
        </div>

        {props.navbar && props.navbar}
      </div>
    </dialog>
  );
};

export default TerminalDialog;
