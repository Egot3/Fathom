<script lang="ts">
  import type { User } from "../lib/contracts/user";

  let { content }: { content: User } = $props();
</script>

<label class="label">
  <span class="label-text">Name</span>
  <input class="input" type="text" value={content.nickname} disabled />
</label>

<label class="label">
  <span class="label-text">UUID</span>
  <input class="input" type="text" value={content.uuid} disabled />
</label>

<label class="label">
  <span class="label-text">Is teacher</span>
  <input class="input" type="text" value={content.is_teacher} disabled />
</label>

<label class="label">
  <span class="label-text">Created at</span>
  <input
    class="input"
    type="date"
    value={content.created_at.slice(0, 10)}
    disabled
  />
</label>

<label class="label">
  <span class="label-text">Updated at</span>
  <input
    class="input"
    type="date"
    value={content.updated_at.slice(0, 10)}
    disabled
  />
</label>
