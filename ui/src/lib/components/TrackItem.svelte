<script lang="ts">
  type Props = {
    fileName: string;
    onRemove: () => void;
  };

  const { fileName, onRemove }: Props = $props();

  let expanded = $state(false);
</script>

<button
  class="bg-red-400"
  onclick={() => {
    expanded = !expanded;
  }}
>
  <p>{fileName}</p>
</button>

{#if expanded}
  <div class="bg-blue-400">
    <p class="line-clamp-1 text-ellipsis" title={fileName}>
      Filename: {fileName}
    </p>
    <button onclick={() => onRemove()}>Remove</button>
  </div>
{/if}
