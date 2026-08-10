<script lang="ts">
  import BentoContainer from "../components/BentoContainer.svelte";
  import Groups from "../components/Groups.svelte";
  import Quizzes from "../components/Quizzes.svelte";
  import Registered from "../components/Registered.svelte";
  import TeacherCurrentlyRunning from "../components/TeacherCurrentlyRunning.svelte";
  import TestingResults from "../components/TestingResults.svelte";
  import Tests from "../components/Tests.svelte";
</script>

<BentoContainer>
  <div
    id="current"
    class="row-start-1 row-end-8 col-start-1 col-end-5 flex flex-col bg-surface-700-300 rounded-2xl p-8"
  >
    <h2 class="text-2xl font-bold">Now running</h2>
    <hr class="hr mt-5 mb-5" />
    <TeacherCurrentlyRunning />
  </div>
  <div
    id="testing_results"
    class="row-start-8 row-end-13 col-start-1 col-end-9 bg-surface-700-300 rounded-2xl p-8 flex flex-col"
  >
    <h2 class="self-center text-2xl font-bold">Testing results</h2>
    <hr class="hr mt-5 mb-5" />
    <TestingResults />
  </div>
  <div
    id="registered"
    class="row-start-1 row-end-8 col-start-5 col-end-7 flex flex-col bg-surface-700-300 rounded-2xl p-8"
  >
    <h2 class="text-2xl font-bold">Registered</h2>
    <hr class="hr mt-5 mb-5" />
    <Registered />
  </div>
  <div
    id="groups"
    class="row-start-1 row-end-8 col-start-7 col-end-9 flex flex-col bg-surface-700-300 rounded-2xl p-8"
  >
    <h2 class="text-2xl font-bold">Groups</h2>
    <hr class="hr mt-5 mb-5" />
    <Groups />
  </div>
  <div
    id="tests"
    class="row-start-1 row-end-7 col-start-9 col-end-13 flex flex-col bg-surface-700-300 rounded-2xl p-8"
  >
    <h2 class="text-2xl font-bold">Tests</h2>
    <hr class="hr mt-5 mb-5" />
    <Tests />
  </div>
  <div
    id="quizzes"
    class="row-start-7 row-end-13 col-start-9 col-end-13 flex flex-col bg-surface-700-300 rounded-2xl p-8"
  >
    <h2 class="text-2xl font-bold">Quizzes</h2>
    <hr class="hr mt-5 mb-5" />
    <Quizzes />
  </div>
</BentoContainer>
