<script>
	import { drawerToggled } from '$lib/store';
	import ChatContent from './chat-content.svelte';
</script>

<main class={$drawerToggled === false ? 'sidenav-hidden' : ''}>
	<div class="window">
		<div class="content">
			<ChatContent />
		</div>
	</div>
</main>

<style lang="postcss">
	main {
		@apply inline;
	}

	.window {
		@apply h-full flex justify-center;
	}

	.content {
		@apply flex flex-col grow shrink bg-blue-50 overflow-x-hidden
  transition-[margin] duration-[400ms];
		@apply sm:me-14 sm:rounded-3xl sm:mb-6 sm:ml-1 mb-0 ml-0;
	}

	.sidenav-hidden .content {
		@apply sm:ms-14;
	}
</style>
