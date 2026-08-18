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

    let title: string = $state("")
    let body: string = $state("")

    let options: QuizOptions | null = $state(null)
    let answer: QuizAnswer | null = $state(null)

    let time: number;
    $effect(()=>{
      loading = true
      const quizUUID: string = UUID;

        time = setTimeout(async () => {
          statusMessage = await FetchParsedQuiz(quizUUID)
          .map((r)=>{
            return r
          })
          .match((r)=>{
            title = r.title
            body = r.body
            return ""
          }, (e)=>e.error);

          loading = false
        }, 1000);
    })
</script>

<h2>{title}</h2>

<p>{body}</p>
