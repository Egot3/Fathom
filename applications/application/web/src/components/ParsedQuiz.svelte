<script lang="ts">
    import { FetchParsedQuiz, type QuizAnswer, type QuizOptions } from "../lib/contracts/quiz";

    const {
      UUID,
      answerValue
    }: {
      UUID: string,
      answerValue: string
    }  = $props()

    let loading: boolean = $state(false)
    let statusMessage: string = $state("")

    let options: QuizOptions | null = $state(null)
    let answer: QuizAnswer | null = $state(null)

    let title: string = $state("")
    let body: string = $state("")

    let time: number;
    $effect(()=>{
      console.log("looping")
      loading = true
      const quizUUID: string = UUID;

        time = setTimeout(async () => {
          statusMessage = await FetchParsedQuiz(quizUUID)
          .match((r)=>{
            console.log(r)

            console.log(title, body)
            title = r.title
            body = r.body

            console.log(title, body)
            return ""
          }, (e)=>e.error);

          loading = false
        }, 1000);
    })
</script>

<h2 class="bg-red-900">{title}</h2>

<p>{body}</p>
