import "@testing-library/jest-dom/vitest";
import { cleanup } from "@solidjs/testing-library";
import { afterAll, afterEach, beforeAll } from "vitest";
import { setupServer } from "msw/node";

export const server = setupServer();

HTMLDialogElement.prototype.showModal = function () {
  this.setAttribute("open", "");
};

HTMLDialogElement.prototype.close = function () {
  this.removeAttribute("open");
};

beforeAll(() => {
  server.listen({ onUnhandledRequest: "error" });
});

afterEach(() => {
  cleanup();
  server.resetHandlers();
});

afterAll(() => {
  server.close();
});

const originalWarn = console.warn;
console.warn = (...args) => {
  if (String(args[0]).includes("computations created outside")) {
    console.trace("Solid warn trace");
  }
  originalWarn(...args);
};
