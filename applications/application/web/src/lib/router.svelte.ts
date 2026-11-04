let path = $state(window.location.pathname);

window.addEventListener("popstate", () => (path = window.location.pathname));

export const CurrentPath = () => path;

export function Navigate(to: string, replace = false) {
  history[replace ? "replaceState" : "pushState"](null, "", to);
  path = to;
}
