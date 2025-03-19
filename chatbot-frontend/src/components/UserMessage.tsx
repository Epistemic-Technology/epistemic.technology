import { Component, JSX } from "solid-js";
import { UserMessageProps } from "../types";
import TerminalMessageWrapper from "./TerminalMessageWrapper";

const UserMessage: Component<UserMessageProps> = (props) => {
  return (
    <TerminalMessageWrapper>
      <div>
        {props.content.split("\n").map((paragraph, index) => (
          <p>{index === 0 ? `> ${paragraph}` : paragraph || "\u00A0"}</p>
        ))}
      </div>
    </TerminalMessageWrapper>
  );
};

export default UserMessage;
