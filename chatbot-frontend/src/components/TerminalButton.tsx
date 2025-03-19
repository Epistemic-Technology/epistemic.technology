import { Component } from 'solid-js';
import terminalIcon from '../assets/svg/terminal.svg';
interface TerminalButtonProps {
  onClick: () => void;
}

const TerminalButton: Component<TerminalButtonProps> = (props) => {
  return (
    <button
      type="button"
      onClick={props.onClick}
      class="fixed bottom-4 right-4 w-12 h-12 rounded-full flex items-center justify-center border-2 hover:bg-opacity-100 transition-all duration-300 focus:outline-none focus:ring-2 focus:ring-opacity-50 z-50"
      aria-label="Open Terminal"
    >
      <img src={terminalIcon} alt="Terminal" class="w-6 h-6" />
    </button>
  );
};

export default TerminalButton; 