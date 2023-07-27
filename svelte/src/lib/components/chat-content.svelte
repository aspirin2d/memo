<script>
	import InputArea from './input-area.svelte';
</script>

<div class="scroll-container">
	<div class="chat-history-container">
		<div class="chat-history-content">Hello, World</div>
	</div>
</div>
<div class="action-bar">
	<InputArea />
</div>

<style lang="postcss">
	.scroll-container {
		@apply flex flex-col overflow-hidden;
	}

	.chat-history-container {
		@apply items-center flex flex-col px-4 overflow-y-auto w-full;
	}

	.chat-history-content {
		@apply m-0 w-full;
		margin-top: 16px;
	}

	.action-bar {
		@apply items-center w-full mt-auto;
	}
</style>
