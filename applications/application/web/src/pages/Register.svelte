<script lang="ts">
  import { ClassForStatus, InputStatus } from "../lib/statuses/input";
  import { CheckPassword } from "../lib/passwordUtils/checkPassword";
  import { Popover, usePopover } from "@skeletonlabs/skeleton-svelte";
  import { FetchRegister } from "../lib/contracts/user";
  import { SetUser } from "../lib/bgdata/user.svelte";
  import { LoaderCircle } from "@lucide/svelte";
  import ErrorPopover from "../components/ErrorPopover.svelte";
  import { Navigate } from "../lib/router.svelte";
  import { onMount } from "svelte";
  import {
    GetPasswordPolicy,
    LoadPasswordPolicy,
  } from "../lib/passwordUtils/policy.svelte";
  import UXInput from "../components/UXInput.svelte";

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

  onMount(() => {
    LoadPasswordPolicy();
  });

  const policy = $derived(GetPasswordPolicy());
  const violations = $derived(policy ? CheckPassword(password, policy) : []);

  $effect(() => {
    passwordMessage = violations.join("\n");
  });

  async function SubmitRegister(e: Event) {
    e.preventDefault();

    registering = true;

    (async () => {
      await FetchRegister(login, password)
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
      registering = false;
    })();

    return null as never;
  }
</script>

<div class="h-full w-full items-center justify-center flex text-surface-50-950">
  <div
    class="h-4/5 w-1/2 bg-surface-950-50 rounded-2xl p-2.5 overflow-scroll flex items-center"
  >
    <form
      onsubmit={SubmitRegister}
      class="w-full flex flex-col h-full justify-center"
    >
      <h1 class="self-center text-6xl justify-self-start m-4 font-bold">
        Register
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

        <label class="label">
          <span class="label-text">Password</span>
          <ErrorPopover open={passwordError} message={passwordMessage}>
            <input
              id="password"
              bind:value={password}
              class={"input border-2 " + passwordClass}
              type="password"
              placeholder="supersecurePaSSwoRD!"
              onfocus={(e: FocusEvent) => {
                passwordError = false;
                passwordState = InputStatus.Idle;
                return;
              }}
              onblur={() => {
                if (password.length === 0) {
                  passwordState = InputStatus.Idle;
                  return;
                }
                if (violations.length > 0) {
                  passwordState = InputStatus.Punish;
                  return;
                }
                passwordState = InputStatus.Treat;
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
              placeholder="supersecurePaSSwoRD12223!"
              onfocus={(e: FocusEvent) => {
                reapeatedError = false;
                repeatedState = InputStatus.Idle;
                return;
              }}
              onblur={() => {
                if (repeated != password) {
                  repeatedState = InputStatus.Punish;
                  reapeatedError = true;
                  repeatedMessage = "Passwords don't match";
                  return;
                }
                repeatedState = passwordState;
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

        <button
          class="btn btn-lg w-2/5 leading-[0.75] preset-outlined-primary-500 self-center"
          onclick={(e) => {
            e.preventDefault();
            Navigate("/login?login=" + login);
          }}
        >
          Log in instead
        </button>
      </div>
    </form>
  </div>
</div>
