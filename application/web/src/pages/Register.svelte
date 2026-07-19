<script lang="ts">
	import { ClassForStatus, InputStatus } from "../lib/statuses/input";
	import { CheckPassword } from "../lib/passwordUtils/checkPassword";
	import { Popover, Portal, usePopover } from "@skeletonlabs/skeleton-svelte";
	import type {
		RegisterRequest,
		RegisterResponse,
	} from "../lib/contracts/user";
	import type { JSONError } from "../lib/statuses/jsonerror";
	import { SetUser } from "../lib/bgdata/user.svelte";
	import { TokenizedFetch } from "../lib/contracts/tokenizedFetch";

	const passwordPopover = usePopover({ id: "password" });
	const loginPopover = usePopover({ id: "login" });
	const repeatedPopover = usePopover({ id: "repeated" });
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

	async function SubmitLogin(e: Event) {
		e.preventDefault();

		const body: RegisterRequest = {
			nickname: login,
			password: password,
		};

		try {
			const response = await TokenizedFetch(
				"http://" + import.meta.env.VITE_DOMAIN + "/api/v1/user/register",
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
				return false;
			}

			const userData = (await response.json()) as RegisterResponse;

			SetUser({
				UUID: userData.user.uuid,
				isTeacher: userData.user.is_teacher,
				nickname: userData.user.nickname,
			});
			window.location.href = "/home";
		} catch (err) {
			console.log("Couldn't fetch register for user: ", err);
			if (err instanceof Error) {
				statusMessage = "couldn't send register because of in-browser error";
			}

			statusMessage = "couldn't send register because of unknown error";
			mainPopover().setOpen(true);
		}
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
									repeatedState = InputStatus.Idle;
									return;
								}}
								onkeyup={() => {
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

				<label class="label">
					<span class="label-text">Password repeated</span>
					<Popover.Provider value={repeatedPopover}>
						<Popover.Anchor>
							<input
								id="password_repeated"
								bind:value={repeated}
								class={"input border-2 " + repeatedClass}
								type="password"
								placeholder="supersecurePaSSwoRD!"
								onblur={(e: FocusEvent) => {
									if (password != repeated) {
										repeatedMessage = "passwords are not equal!";
										repeatedState = InputStatus.Punish;
										repeatedPopover().setOpen(true);
										return;
									}

									repeatedState = InputStatus.Treat;
								}}
								onfocus={(e: FocusEvent) => {
									repeatedPopover().setOpen(false);
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
						</Popover.Anchor>

						<Popover.Positioner>
							<Popover.Content
								class="bg-error-50-950 p-2 rounded-[4px] text-surface-950-50 z-1"
							>
								<Popover.Title tabindex={-1}>{repeatedMessage}</Popover.Title>
							</Popover.Content>
						</Popover.Positioner>
					</Popover.Provider>
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
							Sign up
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
					href={"login?login=" + login} // no, I didn't forget the password, form vals stay in history
				>
					Log in instead
				</a>
			</div>
		</form>
	</div>
</div>
