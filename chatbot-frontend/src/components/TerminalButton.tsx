import { Component } from "solid-js";
import terminalIcon from "../assets/svg/terminal.svg";
interface TerminalButtonProps {
  onClick: () => void;
}

const TerminalButton: Component<TerminalButtonProps> = (props) => {
  return (
    <div
      class="flex items-center gap-2 fixed bottom-4 right-4"
      role="complementary"
      aria-label="Chat shortcuts"
    >
      <kbd
        class="px-2 py-1.5 text-sm font-semibold bg-gray-100 border border-gray-300 rounded-md shadow-sm text-gray-800 dark:bg-gray-800 dark:text-gray-200 dark:border-gray-600"
        aria-label="Press slash key to open chat"
        role="tooltip"
      >
        /
      </kbd>
      <span class="text-2xl" aria-hidden="true">
        ⇒
      </span>
      <button
        type="button"
        onClick={props.onClick}
        class="w-12 h-12 rounded-full flex items-center justify-center border-2 hover:bg-opacity-100 transition-all duration-300 focus:outline-none focus:ring-2 focus:ring-opacity-50 z-50 bg-white dark:bg-gray-800 border-gray-300 dark:border-gray-600 cursor-pointer"
        aria-label="Open Chat (press slash key as shortcut)"
      >
        <img
          src={terminalIcon}
          alt="Terminal"
          class="w-6 h-6 terminal-icon"
          style={{
            filter: "invert(1)",
          }}
        />
      </button>
    </div>
  );
};

export default TerminalButton;
