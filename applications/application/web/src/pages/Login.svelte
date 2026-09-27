<script lang="ts">
  import { ClassForStatus, InputStatus } from "../lib/statuses/input";
  import { CheckPassword } from "../lib/passwordUtils/checkPassword";
  import { Popover, Portal, usePopover } from "@skeletonlabs/skeleton-svelte";
  import { FetchLogin } from "../lib/contracts/user";
  import { IsJSONError, type JSONError } from "../lib/statuses/jsonerror";
  import { SetUser } from "../lib/bgdata/user.svelte";
  import { LoaderCircle } from "@lucide/svelte";
  import ErrorPopover from "../components/ErrorPopover.svelte";

  const mainPopover = usePopover({ id: "main" });
  let loginError = $state(false);
  let passwordError = $state(false);

  const params = new URLSearchParams(window.location.search);
  let login = $state(params.get("login") ?? "");
  let password = $state("");

  let statusMessage = $state("");
  let passwordMessage = $state("");
  let loginMessage = $state("");

  let loggingIn = $state(false);

  let passwordState = $state(InputStatus.Idle);
  let loginState = $derived(
    login == ""
      ? InputStatus.Idle
      : login.length < 3
        ? InputStatus.Punish
        : InputStatus.Treat,
  );

  let passwordClass = $derived(ClassForStatus(passwordState));
  let loginClass = $derived(ClassForStatus(loginState));

  async function SubmitLogin(e: Event) {
    e.preventDefault();

    loggingIn = true;

    FetchLogin(login, password)
      .then((r) => {
        if (IsJSONError(r)) {
          statusMessage = r.error;
          mainPopover().setOpen(true);
          return;
        }

        SetUser({
          UUID: r.user.uuid,
          isTeacher: r.user.is_teacher,
          nickname: r.user.nickname,
        });
        window.location.href = "/home";
      })
      .finally(() => (loggingIn = false));

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
      <div class="flex flex-col msx-w-md space-y-4 mx-auto w-full">
        <label class="label">
          <span class="label-text">Login</span>
          <ErrorPopover message={loginMessage} open={loginError}>
            <input
              class={"input border-2 " + loginClass}
              type="text"
              placeholder="myoryourlogin"
              bind:value={login}
              onblur={(e: FocusEvent) => {
                if (login.length < 3) {
                  loginState = InputStatus.Punish;
                  loginMessage = "can't have length of login less than 3";
                  loginError = true;
                  return;
                }
                loginState = InputStatus.Treat;
              }}
              onfocus={(e: FocusEvent) => {
                loginError = false;
                loginState = InputStatus.Idle;
                return;
              }}
              onkeydown={() => {
                if (login.length >= 3) {
                  loginState = InputStatus.Treat;
                }
              }}
            />
          </ErrorPopover>
        </label>
        <label class="label">
          <span class="label-text">Password</span>
          <ErrorPopover open={passwordError} message={passwordMessage}>
            <input
              id="password"
              bind:value={password}
              class={"input border-2 " + passwordClass}
              type="password"
              placeholder="supersecurePaSSwoRD!"
              onblur={(e: FocusEvent) => {
                const err = CheckPassword(password);
                if (err !== "") {
                  passwordMessage = err;
                  passwordState = InputStatus.Punish;
                  passwordError = true;
                  return;
                }

                passwordState = InputStatus.Treat;
              }}
              onfocus={(e: FocusEvent) => {
                passwordError = false;
                passwordState = InputStatus.Idle;
                return;
              }}
              onkeydown={() => {
                const err = CheckPassword(password);
                if (err == "") {
                  passwordState = InputStatus.Treat;
                }
              }}
            />
          </ErrorPopover>
        </label>

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
              disabled={!(
                passwordState == InputStatus.Treat &&
                loginState == InputStatus.Treat
              )}
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

      <a
        class="btn btn-lg w-2/5 leading-[0.75] text-xl preset-outlined-primary-500 self-center"
        href={"register?login=" + login}
      >
        Register
      </a>
    </form>
  </div>
</div>
<!-- login via github is forbidden by law -->
