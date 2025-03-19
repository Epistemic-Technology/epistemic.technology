import { Component } from "solid-js";
import { MessageWrapperProps } from "../types";

const TerminalMessageWrapper: Component<MessageWrapperProps> = (props) => {
  return (
    <div class={`terminal-message ${props.className || ""}`}>
      <div class="flex">
        <div class="text-2xl">{props.children}</div>
      </div>
    </div>
  );
};

export default TerminalMessageWrapper;
