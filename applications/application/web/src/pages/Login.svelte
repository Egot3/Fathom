<script lang="ts">
  import { ClassForStatus, InputStatus } from "../lib/statuses/input";
  import { Popover, Portal, usePopover } from "@skeletonlabs/skeleton-svelte";
  import { FetchLogin } from "../lib/contracts/user";
  import { SetUser } from "../lib/bgdata/user.svelte";
  import { LoaderCircle } from "@lucide/svelte";
  import { Navigate } from "../lib/router.svelte";
  import UXInput from "../components/UXInput.svelte";

  const mainPopover = usePopover({ id: "main" });

  const params = new URLSearchParams(window.location.search);
  let login = $state(params.get("login") ?? "");
  let password = $state("");

  let statusMessage = $state("");
  let loginMessage = $state("");

  let loggingIn = $state(false);

  let loginState = $derived(
    login == ""
      ? InputStatus.Idle
      : login.length < 3
        ? InputStatus.Punish
        : InputStatus.Treat,
  );

  let loginClass = $derived(ClassForStatus(loginState));

  async function SubmitLogin(e: Event) {
    e.preventDefault();

    loggingIn = true;

    (async () => {
      await FetchLogin(login, password)
        .andTee((r) => {
          SetUser({
            UUID: r.user.uuid,
            isTeacher: r.user.is_teacher,
            nickname: r.user.nickname,
          });
          Navigate("/home");
        })
        .orTee((e) => {
          statusMessage = e.error;
          mainPopover().setOpen(true);
          return;
        });
      loggingIn = false;
    })();

    return null as never;
  }
</script>

<div class="h-full w-full items-center justify-center flex text-surface-50-950">
  <div
    class="h-4/5 w-1/2 bg-surface-950-50 rounded-2xl p-2.5 overflow-auto flex items-center"
  >
    <form
      onsubmit={SubmitLogin}
      class="w-full flex flex-col h-full justify-center"
    >
      <h1 class="self-center text-6xl justify-self-start m-4 font-bold">
        Log in
      </h1>
      <div class="flex flex-col space-y-4 mx-auto w-full">
        <UXInput
          label="Login"
          bind:value={login}
          type="text"
          bind:message={loginMessage}
          checker={(v: string) => {
            if (login.length >= 3) {
              return true;
            }

            if (login.length === 0) {
              return false;
            }
            if (login.length < 3) {
              loginMessage = "can't have length of login less than 3";
              return false;
            }

            return true;
          }}
          placeholder="myoryourlogin"
        ></UXInput>

        <UXInput
          label="Password"
          bind:value={password}
          checker={(v: string) => password.length > 0}
          type="password"
          placeholder="supersecurePaSSwoRD12312!"
        ></UXInput>

        <Popover>
          <Popover.Trigger
            class="btn preset-filled self-start text-xs text-primary-500 -mt-3.5 h-fit p-0.75"
            >Forgot password?</Popover.Trigger
          >
          <Portal>
            <Popover.Positioner>
              <Popover.Content
                class="card w-96 p-4 bg-surface-100-900 shadow-xl"
              >
                <Popover.Title
                  >Nicely ask your teacher/hoster/provicer to kindly confirm
                  your identity and awesomly forge you another one. If YOU are
                  the hoster/provider address the docs</Popover.Title
                >
                <!-- TODO, docs link -->
                <Popover.Arrow
                  class="[--arrow-size:--spacing(2)] [--arrow-background:var(--color-surface-100-900)]"
                >
                  <Popover.ArrowTip />
                </Popover.Arrow>
              </Popover.Content>
            </Popover.Positioner>
          </Portal>
        </Popover>

        <Popover.Provider value={mainPopover}>
          <Popover.Anchor class="self-center w-full flex justify-center gap-0">
            <button
              type="submit"
              class="btn btn-lg w-4/5 preset-filled-primary-500"
              disabled={!(loginState == InputStatus.Treat)}
            >
              {#if loggingIn}
                <LoaderCircle class="animate-spin" />
              {:else}
                Sign in
              {/if}
            </button>
          </Popover.Anchor>

          <Popover.Positioner>
            <Popover.Content
              class="bg-error-50-950 p-2 rounded-[4px] text-surface-950-50 z-1"
            >
              <Popover.Title tabindex={-1}>{statusMessage}</Popover.Title>
            </Popover.Content>
          </Popover.Positioner>
        </Popover.Provider>
      </div>

      <button
        class="btn btn-lg w-2/5 leading-[0.75] text-xl preset-outlined-primary-500 self-center"
        onclick={(e) => {
          e.preventDefault();
          Navigate("/register?login=" + login);
        }}
      >
        Register
      </button>
    </form>
  </div>
</div>
<!-- login via github is forbidden by law -->
