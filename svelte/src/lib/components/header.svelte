<script>
	import { drawerToggled } from '$lib/store';
</script>

<div class="header">
	<!-- header left -->
	<div class="start">
		<button on:click={() => ($drawerToggled = !$drawerToggled)}>
			<svg
				xmlns="http://www.w3.org/2000/svg"
				fill="none"
				viewBox="0 0 24 24"
				stroke-width="1.5"
				stroke="currentColor"
				class="w-6 h-6"
			>
				<path
					stroke-linecap="round"
					stroke-linejoin="round"
					d="M3.75 6.75h16.5M3.75 12h16.5m-16.5 5.25h16.5"
				/>
			</svg>
		</button>
		<div class="ml-4 font-semibold text-neutral-600 cursor-default text-lg">ICE</div>
	</div>
	<!-- header center -->
	<div />
	<!-- header right -->
	<div />
</div>

<style lang="postcss">
	.header {
		@apply fixed w-full h-16 flex flex-row items-center justify-between px-2 text-neutral-600;
	}

	.start {
		@apply flex flex-row items-center;
	}
</style>
