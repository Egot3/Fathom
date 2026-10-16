export function CreateFitRows(opts: {
  rowHeight: number;
  reserve?: number;
  min?: number;
  max?: number;
}) {
  const { rowHeight, reserve = 0, min = 1, max = 100 } = opts;
  let height = $state(0);

  const pageSize = $derived(
    Math.min(max, Math.max(min,
      Math.floor(height / rowHeight) - reserve)),
  );

  function measure(node: HTMLElement) {
    const ro = new ResizeObserver(([e]) => {
      height = e.contentRect.height;
    });
    ro.observe(node);
    return { destroy: () => ro.disconnect() };
  }

  return {
    measure,
    get pageSize() { return pageSize; },
    get ready() { return height > 0; },
  };
}
