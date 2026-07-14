<script lang="ts">
  import { onMount } from "svelte";

  const regex = /^http.\/\/.+?\/(.*)$/gm

  let currentURL = $state('')

  let slug = $derived(
    regex.exec(currentURL)?.[1] ?? ""
  )

  onMount(()=>{
    currentURL = window.location.href
  })
</script>

<p class="bg-error-100">current url: {currentURL}</p>
<p>slug: {slug}</p>
<p style="background: #fbbf24;">a bunch of text</p>