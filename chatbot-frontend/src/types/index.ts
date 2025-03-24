import { Component, JSX } from "solid-js";

export interface MessageWrapperProps {
  children: JSX.Element;
  className?: string;
}

export interface TerminalMessageProps {
  type: "user" | "bot";
  content: string;
  sources?: Source[];
}

export interface Source {
  index?: number;
  Title: string;
  FilePath: string;
  Author: string;
  ID: number;
  Content: string;
  PublicationDate: string;
}

export interface ChatRequest {
  query: string;
  history: string;
}

export interface BotMessageProps extends TerminalMessageProps {}

export interface UserMessageProps extends TerminalMessageProps {}

export interface UserInputProps {
  onSubmit: (query: string) => Promise<void>;
  isLoading: boolean;
  autoFocus?: boolean;
}

export interface Buffer {
  id: string;
  name: string;
  render: () => JSX.Element;
  onActivate?: () => void;
  onDeactivate?: () => void;
}

export interface ChatBufferProps {
  chatHistory?: TerminalMessageProps[];
  isLoading?: boolean;
  apiUrl?: string;
  onExit?: () => void;
}
