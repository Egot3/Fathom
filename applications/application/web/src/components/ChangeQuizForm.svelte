<script lang="ts">

	import { Dialog } from "@skeletonlabs/skeleton-svelte";
	import {
	FetchQuiz,
	FetchQuizPatch,
	FetchQuizPut,
	type Meta,
	} from "../lib/contracts/quiz";
	import { IsJSONError } from "../lib/statuses/jsonerror";
	import UXInput from "./UXInput.svelte";
    import _ from "lodash"

	const {
		callback,
		UUID,
        name
	}: {
		callback: () => void;

        name: string;
		UUID: string;
	} = $props();

	const pathRegex = /^([a-zA-Z0-9_\-]+)*\/?[a-zA-Z0-9_\-]+$/;

    

    let statusMessage = $state("")
    let nameMessage: string = $state("");

	let nameReady: boolean = $derived(pathRegex.test(name));

    let newName = $derived(name)
    let body = $state("")
	let allOrNone: boolean = $state(false);
	let randomized: boolean = $state(true);
	let score: number = $state(1);

    let loading: boolean = $state(true)
    let error: string = $state("")

    const response = $derived(FetchQuiz(UUID).
        then((r)=>{
            if(IsJSONError(r)) {
                error = r.error
                return null as never
            }
            body = r.body
            allOrNone = r.meta.all_or_none
            randomized = r.meta.randomized
            score = r.meta.score
            return r
        }).finally(()=>loading = false))
    
</script>

{#if loading}
    <div class="animate-pulse h-full w-full bg-surface-400-600 rounded-xl">

    </div>
{:else} 
    {#if error != ""}
        <p>{error}</p>
    {:else}
    <form onsubmit={async (e) =>{
        e.preventDefault()

        const quizResponse = await response 

        const newMeta: Meta = {
            all_or_none: allOrNone,
            score: score,
            randomized: randomized,
            kind: quizResponse.meta.kind // т.к. раньше надо было думать
        }

        if (newName!=name || score != quizResponse.meta.score) {
            const patchResponse = await FetchQuizPatch(UUID, 
                score!=quizResponse.meta.score ? score : undefined, 
                newName!=name ? newName : undefined)
            if(patchResponse!=null) {
                statusMessage = patchResponse.error
                return
            }
        }

        
        if (
            body!=quizResponse.body ||
            !_.isEqual(quizResponse.meta, newMeta)
        ) {
            const putResponse = await FetchQuizPut(UUID, newMeta, body)
            if (putResponse!=null) {
                statusMessage = putResponse.error
                return
            }
        }
    }} class="w-full flex flex-col h-full space-y-2">
	
    <UXInput
		label="Name/path"
		placeholder="color_theory/sky.md"
		bind:value={newName}
		bind:ready={nameReady}
		message={nameMessage}
		checker={(v: string) => {
			nameMessage =
				"name must have no special symbols and doesn't contain any extensions";
			return pathRegex.test(v);
		}}
	/>
    


	<footer class="flex justify-end gap-2 w-full">
		{#if statusMessage != ""}
			<span class="justify-self-start">{statusMessage}</span>
		{/if}
		<Dialog.CloseTrigger class="btn preset-tonal">Cancel</Dialog.CloseTrigger>
		<button type="submit" class="btn preset-filled">Create</button>
	</footer>
</form>

    {/if}
{/if}
