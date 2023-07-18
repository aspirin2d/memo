<script>
</script>

<main>
	<div class="window">
		<div class="content">Hello, World</div>
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
		@apply flex flex-col grow shrink bg-red-200 overflow-x-hidden;
		margin-block: 1px 24px;
		margin-inline: 0 56px;
		border-radius: 24px;
	}
</style>
