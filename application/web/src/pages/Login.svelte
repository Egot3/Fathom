<script lang="ts">
	import { ClassForStatus, InputStatus } from "../lib/statuses/input";
	import { CheckPassword } from "../lib/passwordUtils/checkPassword";
	import { Popover, Portal, usePopover } from "@skeletonlabs/skeleton-svelte";

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

	async function SubmitLogin(
		event: SubmitEvent & { currentTarget: HTMLFormElement },
	) {}
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
			</div>
		</form>
	</div>
</div>
