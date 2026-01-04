import "@testing-library/jest-dom/vitest";
import { cleanup } from "@testing-library/react";
import { afterEach, beforeEach, vi } from "vitest";
import {
  createMockApiState,
  createMockFetch,
  type MockApiState,
} from "./mockApi";

vi.mock("~/context/glossary", () => ({
  GlossaryProvider: ({ children }: { children: React.ReactNode }) => children,
}));

HTMLDialogElement.prototype.showModal = function () {
  this.setAttribute("open", "");
};

HTMLDialogElement.prototype.close = function () {
  this.removeAttribute("open");
};

let mockApiState: MockApiState;

beforeEach(() => {
  mockApiState = createMockApiState();
  vi.stubGlobal("fetch", createMockFetch(mockApiState));
});

afterEach(() => {
  cleanup();
  vi.unstubAllGlobals();
});

export function getMockApiState(): MockApiState {
  return mockApiState;
}
