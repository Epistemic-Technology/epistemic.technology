import { Component, For } from "solid-js";

export interface NavbarCommand {
  name: string;
  command: string;
}

export interface NavbarProps {
  commands: NavbarCommand[];
  onCommandExecute: (command: string) => void;
}

const Navbar: Component<NavbarProps> = (props) => {
  return (
    <div class="bg-blue-300 text-black p-2 flex justify-center space-x-6 w-full mt-auto chat-navbar">
      <For each={props.commands}>
        {(cmd) => (
          <button
            class="font-bold hover:underline focus:outline-none cursor-pointer"
            onClick={() => props.onCommandExecute(cmd.command)}
          >
            {cmd.name}
          </button>
        )}
      </For>
    </div>
  );
};

export default Navbar;
