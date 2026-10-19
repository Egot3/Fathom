<script lang="ts">
  import { ArrowDown, ArrowUp } from "@lucide/svelte";
  import {
    Kind,
    OptionAccordance,
    OptionCheck,
    OptionOrder,
    OptionRadio,
    type QuizAnswer,
    type QuizOptions,
  } from "../lib/contracts/quiz";
  import { ClassForStatus, InputStatus } from "../lib/statuses/input";

  let {
    kind,
    options,
    value = $bindable(),
    correct,
    readonly = false,
    onanswer,
  }: {
    kind: Kind;
    options: QuizOptions;
    value?: QuizAnswer;
    correct?: QuizAnswer;
    readonly?: boolean;
    onanswer?: (a: QuizAnswer) => void;
  } = $props();

  const group = $props.id();
  const review = $derived(readonly && correct !== undefined);

  const v = $derived(value as any);
  const c = $derived(correct as any);

  function set(next: QuizAnswer) {
    value = next;
    onanswer?.(next);
  }

  const mark = (chosen: boolean, right: boolean) =>
    !review ? "" : right ? "text-success-500" : chosen ? "text-error-500" : "";

  const picked = $derived(new Set<number>(v?.check?.chosen ?? []));
  const rights = $derived(new Set<number>(c?.check?.chosen ?? []));
  function toggle(id: number) {
    const next = new Set(picked);
    next.has(id) ? next.delete(id) : next.add(id);
    set({ check: { chosen: [...next].sort((a, b) => a - b) } });
  }

  const items = $derived(kind === Kind.Order ? OptionOrder(options) : []);
  const order = $derived<number[]>(
    v?.order?.item_indexes ?? items.map((_, i) => i),
  );
  function move(pos: number, d: -1 | 1) {
    const next = [...order];
    [next[pos], next[pos + d]] = [next[pos + d], next[pos]];
    set({ order: { item_indexes: next } });
  }

  const acc = $derived(
    kind === Kind.Accordance
      ? OptionAccordance(options)
      : { static: [] as string[], dynamic: [] as string[] },
  );
  const accorded = $derived<number[]>(
    v?.accordance?.accorded ?? acc.static.map(() => -1),
  );
  function accord(i: number, j: number) {
    const next = [...accorded];
    next[i] = j;
    set({ accordance: { accorded: next } });
  }
</script>

{#if kind === Kind.Input}
  <input
    class={"input border-2 " +
      (review
        ? ClassForStatus(
            v?.input?.input === c?.input?.input
              ? InputStatus.Treat
              : InputStatus.Punish,
          )
        : "")}
    type="text"
    disabled={readonly}
    value={v?.input?.input ?? ""}
    oninput={(e) => set({ input: { input: e.currentTarget.value } })}
  />
  {#if review && v?.input?.input !== c?.input?.input}
    <p class="text-success-500">{c?.input?.input}</p>
  {/if}
{:else if kind === Kind.Radio}
  <div class="space-y-2">
    {#each OptionRadio(options) as o (o.id)}
      <label
        class={"flex items-center space-x-2 " +
          mark(v?.radio?.chosen === o.id, c?.radio?.chosen === o.id)}
      >
        <input
          class="radio"
          type="radio"
          name={group}
          disabled={readonly}
          checked={v?.radio?.chosen === o.id}
          onchange={() => set({ radio: { chosen: o.id } })}
        />
        <p>{o.label}</p>
      </label>
    {/each}
  </div>
{:else if kind === Kind.Check}
  <div class="space-y-2">
    {#each OptionCheck(options) as o (o.id)}
      <label
        class={"flex items-center space-x-2 " +
          mark(picked.has(o.id), rights.has(o.id))}
      >
        <input
          class="checkbox"
          type="checkbox"
          disabled={readonly}
          checked={picked.has(o.id)}
          onchange={() => toggle(o.id)}
        />
        <p>{o.label}</p>
      </label>
    {/each}
  </div>
{:else if kind === Kind.Order}
  <ol class="space-y-1 list-decimal ml-6">
    {#each order as idx, pos (idx)}
      <li>
        <span class="flex items-center space-x-2">
          <span>{items[idx]}</span>
          {#if !readonly}
            <button
              type="button"
              class="btn-icon"
              disabled={pos === 0}
              onclick={() => move(pos, -1)}><ArrowUp size={14} /></button
            >
            <button
              type="button"
              class="btn-icon"
              disabled={pos === order.length - 1}
              onclick={() => move(pos, 1)}><ArrowDown size={14} /></button
            >
          {/if}
        </span>
      </li>
    {/each}
  </ol>
  {#if review}
    <p class="text-success-500 mt-2">
      {(c?.order?.item_indexes as number[]).map((i) => items[i]).join(" → ")}
    </p>
  {/if}
{:else if kind === Kind.Accordance}
  <div class="space-y-2">
    {#each acc.static as key, i (key)}
      <label
        class={"flex items-center space-x-2 " +
          mark(true, accorded[i] === c?.accordance?.accorded?.[i])}
      >
        <span class="w-1/3">{key}</span>
        <select
          class="select"
          disabled={readonly}
          onchange={(e) => accord(i, Number(e.currentTarget.value))}
        >
          <option value={-1} selected={accorded[i] === -1}>—</option>
          {#each acc.dynamic as d, j}
            <option value={j} selected={accorded[i] === j}>{d}</option>
          {/each}
        </select>
        {#if review && accorded[i] !== c?.accordance?.accorded?.[i]}
          <span class="text-success-500"
            >{acc.dynamic[c.accordance.accorded[i]]}</span
          >
        {/if}
      </label>
    {/each}
  </div>
{:else}
  <div>Unknown quiz kind!</div>
{/if}
