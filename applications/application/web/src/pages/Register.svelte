<script lang="ts">
  import { ClassForStatus, InputStatus } from "../lib/statuses/input";
  import { CheckPassword } from "../lib/passwordUtils/checkPassword";
  import { Popover, Portal, usePopover } from "@skeletonlabs/skeleton-svelte";
  import {
    FetchRegister,
    type RegisterRequest,
    type RegisterResponse,
  } from "../lib/contracts/user";
  import { IsJSONError, type JSONError } from "../lib/statuses/jsonerror";
  import { SetUser } from "../lib/bgdata/user.svelte";
  import { LoaderCircle } from "@lucide/svelte";
  import ErrorPopover from "../components/ErrorPopover.svelte";

  let loginError = $state(false);
  let passwordError = $state(false);
  let reapeatedError = $state(false);
  const mainPopover = usePopover({ id: "main" });

  const params = new URLSearchParams(window.location.search);
  let login = $state(params.get("login") ?? "");
  let password = $state("");
  let repeated = $state("");

  let statusMessage = $state("");
  let passwordMessage = $state("");
  let loginMessage = $state("");
  let repeatedMessage = $state("");

  let loginState = $derived(
    login == ""
      ? InputStatus.Idle
      : login.length < 3
        ? InputStatus.Punish
        : InputStatus.Treat,
  );
  let passwordState = $state(InputStatus.Idle);
  let repeatedState = $state(InputStatus.Idle);

  let passwordClass = $derived(ClassForStatus(passwordState));
  let loginClass = $derived(ClassForStatus(loginState));
  let repeatedClass = $derived(ClassForStatus(repeatedState));

  let registering = $state(false);

  async function SubmitLogin(e: Event) {
    e.preventDefault();
    registering = true;

    const body: RegisterRequest = {
      nickname: login,
      password: password,
    };

    FetchRegister(login, password)
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
      .finally(() => (registering = false));
  }
</script>

<div class="h-full w-full items-center justify-center flex text-surface-50-950">
  <div
    class="h-4/5 w-1/2 bg-surface-950-50 rounded-2xl p-2.5 overflow-scroll flex items-center"
  >
    <form
      onsubmit={SubmitLogin}
      class="w-full flex flex-col h-full justify-center"
    >
      <h1 class="self-center text-6xl justify-self-start m-4 font-bold">
        Register
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
                if (err !== null) {
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

        <label class="label">
          <span class="label-text">Password repeated</span>
          <ErrorPopover open={reapeatedError} message={repeatedMessage}>
            <input
              id="password_repeated"
              bind:value={repeated}
              class={"input border-2 " + repeatedClass}
              type="password"
              placeholder="supersecurePaSSwoRD!"
              onblur={(e: FocusEvent) => {
                if (password !== repeated) {
                  repeatedMessage = "passwords are not equal!";
                  repeatedState = InputStatus.Punish;
                  reapeatedError = true;
                  return;
                }

                repeatedState = InputStatus.Treat;
              }}
              onfocus={(e: FocusEvent) => {
                reapeatedError = false;
                repeatedState = InputStatus.Idle;
                return;
              }}
              onkeyup={() => {
                console.log("keydown");
                if (password == repeated) {
                  repeatedState = InputStatus.Treat;
                }
              }}
            />
          </ErrorPopover>
        </label>

        <Popover.Provider value={mainPopover}>
          <Popover.Anchor class="self-center w-full flex justify-center gap-0">
            <button
              type="submit"
              class="btn btn-lg w-4/5 preset-filled-primary-500 self-center"
              disabled={!(
                passwordState == InputStatus.Treat &&
                loginState == InputStatus.Treat &&
                repeatedState == InputStatus.Treat
              )}
            >
              {#if registering}
                <LoaderCircle class="animate-spin" />
              {:else}
                Sign up
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

        <a
          class="btn btn-lg w-2/5 leading-[0.75] preset-outlined-primary-500 self-center"
          href={"login?login=" + login}
        >
          Log in instead
        </a>
      </div>
    </form>
  </div>
</div>
