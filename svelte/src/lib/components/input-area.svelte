<script>
	import { onMount } from 'svelte';

	var observe;
	var textarea;

	onMount(() => {
		if (window.attachEvent) {
			observe = function (element, event, handler) {
				element.attachEvent('on' + event, handler);
			};
		} else {
			observe = function (element, event, handler) {
				element.addEventListener(event, handler, false);
			};
		}

		function resize() {
			textarea.style.height = 'auto';
			textarea.style.height = textarea.scrollHeight + 'px';
		}
		/* 0-timeout to get the already changed text */
		function delayedResize() {
			window.setTimeout(resize, 0);
		}

		observe(textarea, 'change', resize);
		observe(textarea, 'cut', delayedResize);
		observe(textarea, 'paste', delayedResize);
		observe(textarea, 'drop', delayedResize);
		observe(textarea, 'keydown', delayedResize);

		textarea.focus();
		textarea.select();
		resize();
	});
</script>

<div class="input-area">
	<!-- command line button -->
	<button class="command-button">
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
				d="M6.75 7.5l3 2.25-3 2.25m4.5 0h3m-9 8.25h13.5A2.25 2.25 0 0021 18V6a2.25 2.25 0 00-2.25-2.25H5.25A2.25 2.25 0 003 6v12a2.25 2.25 0 002.25 2.25z"
			/>
		</svg>
	</button>
	<div class="input-form">
		<div class="textarea-container">
			<textarea bind:this={textarea} rows="1" style="height:1em;" />
		</div>
	</div>
	<button class="send-button">
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
				d="M6 12L3.269 3.126A59.768 59.768 0 0121.485 12 59.77 59.77 0 013.27 20.876L5.999 12zm0 0h7.5"
			/>
		</svg>
	</button>
</div>

<style lang="postcss">
	.input-area {
		@apply my-4 mx-10 flex flex-row;
	}

	.command-button {
		@apply ms-0 me-4;
	}

	.send-button {
		@apply ms-4 me-0;
	}

	.input-form {
		@apply bg-white rounded-[28px] flex flex-col flex-1 px-8 box-border border
  border-gray-400;
	}

	.input-form:hover {
		@apply border-gray-600;
	}

	.input-form:focus-within {
		@apply border-blue-400;
	}

	.textarea-container {
		@apply flex-1 min-h-[56px] py-4;
	}

	textarea {
		@apply w-full min-h-[16px] max-h-[192px] border-none inline-block m-0
  outline-none resize-none p-0 border-0 align-middle transition-all;
	}
</style>
