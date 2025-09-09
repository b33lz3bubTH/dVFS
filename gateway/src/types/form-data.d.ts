declare global {
  class FormData {
    append(name: string, value: any, filename?: string): void;
  }
}