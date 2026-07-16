<script lang="ts">
	import { ClassForStatus, InputStatus } from "../lib/statuses/input";
	import { CheckPassword } from "../lib/passwordUtils/checkPassword";
	import { Popover, Portal, usePopover } from "@skeletonlabs/skeleton-svelte";
	import type { LoginRequest, LoginResponse } from "../lib/contracts/user";
	import type { JSONError } from "../lib/statuses/jsonerror";
	import { setUser } from "../lib/contexts/user";

	const passwordPopover = usePopover({ id: "password" });
	const loginPopover = usePopover({ id: "login" });

	let login = $state("");
	let password = $state("");

	let statusMessage = $state("");
	let passwordMessage = $state("");
	let loginMessage = $state("");

	let passwordState = $state(InputStatus.Idle);
	let loginState = $state(InputStatus.Idle);

	let passwordClass = $derived(ClassForStatus(passwordState));
	let loginClass = $derived(ClassForStatus(loginState));

	async function SubmitLogin() {
		const body: LoginRequest = {
			nickname: login,
			password: password,
		};

		try {
			const response = await fetch(
				import.meta.env.DOMAIN + "/api/v1/user/login",
				{
					method: "POST",
					headers: {
						"Content-Type": "application/json",
						Accept: "application/json",
					},
					body: JSON.stringify(body),
				},
			);

			if (!response.ok) {
				statusMessage = ((await response.json()) as JSONError).error; // then they say that front is better
				return;
			}

			const userData = (await response.json()) as LoginResponse;
			setUser(userData.user);
			window.location.href = "/home";
		} catch (err) {
			console.log("Couldn't fetch login for user: ", err);
			if (err instanceof Error) {
				statusMessage = "couldn't send login because of in-browser error";
			}

			statusMessage = "couldn't send login because of unknown error";
		}
	}
</script>

<div
	class="h-full w-full items-center justify-center flex flex-col text-surface-50-950"
>
	<div
		class="h-4/5 w-1/2 bg-surface-950-50 rounded-2xl p-2.5 overflow-scroll flex items-center"
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
					<Popover.Provider value={loginPopover}>
						<Popover.Anchor>
							<input
								class={"input border-2 " + loginClass}
								type="text"
								placeholder="myoryourlogin"
								bind:value={login}
								onblur={(e: FocusEvent) => {
									if (login.length < 3) {
										loginState = InputStatus.Punish;
										loginMessage = "can't have length of login less than 3";
										loginPopover().setOpen(true);
										return;
									}
									loginState = InputStatus.Treat;
								}}
								onfocus={(e: FocusEvent) => {
									loginPopover().setOpen(false);
									loginState = InputStatus.Idle;
									return;
								}}
								onkeydown={() => {
									if (login.length >= 3) {
										loginState = InputStatus.Treat;
									}
								}}
							/>
						</Popover.Anchor>

						<Popover.Positioner>
							<Popover.Content
								class="bg-error-50-950 p-2 rounded-[4px] text-surface-950-50"
							>
								<Popover.Title tabindex={-1}>{loginMessage}</Popover.Title>
							</Popover.Content>
						</Popover.Positioner>
					</Popover.Provider>
				</label>
				<label class="label">
					<span class="label-text">Password</span>
					<Popover.Provider value={passwordPopover}>
						<Popover.Anchor>
							<input
								id="password"
								bind:value={password}
								class={"input border-2 " + passwordClass}
								type="password"
								placeholder="supersecurePaSSwoRD!"
								onblur={(e: FocusEvent) => {
									const err = CheckPassword(password);
									if (err != "") {
										passwordMessage = err;
										passwordState = InputStatus.Punish;
										passwordPopover().setOpen(true);
										return;
									}

									passwordState = InputStatus.Treat;
								}}
								onfocus={(e: FocusEvent) => {
									passwordPopover().setOpen(false);
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
						</Popover.Anchor>

						<Popover.Positioner>
							<Popover.Content
								class="bg-error-50-950 p-2 rounded-[4px] text-surface-950-50 z-1"
							>
								<Popover.Title tabindex={-1}>{passwordMessage}</Popover.Title>
							</Popover.Content>
						</Popover.Positioner>
					</Popover.Provider>
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

				<!-- 				<label class="flex items-center space-x-2 w-fit pr-2">
					chat, is this UX moment?
					<button
						{...toggle.trigger}
						type="button"
						aria-label="toggle stay logged-in"
					>
						<Maximize size="28">
							{#if toggle.value}
								<Check size="18" x="3" y="3" />
							{:else}
								<X size="18" x="3" y="3" />
							{/if}
						</Maximize>
					</button>
					<b>Stay logged-in</b>
				</label> -->

				<button
					type="submit"
					class="btn btn-lg w-4/5 preset-filled-primary-500 self-center"
					disabled={!(
						passwordState == InputStatus.Treat &&
						loginState == InputStatus.Treat
					)}
				>
					Login
				</button>

				<a
					class="btn btn-lg w-2/5 leading-[0.75] preset-outlined-primary-500 self-center"
					href={"register?login=" + login} // no, I didn't forget the password, form vals stay in history
				>
					Register instead
				</a>
			</div>
		</form>
	</div>
</div>
<!-- login via github is forbidden by law -->
