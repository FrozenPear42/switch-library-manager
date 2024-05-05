import * as Progress from "@radix-ui/react-progress";

export default function ProgressBar({ progress }: { progress: number }) {
  return (

    <Progress.Root className="radix-ProgressRoot" value={Math.floor(progress)}>
      <Progress.Indicator
        className="radix-ProgressIndicator"
        style={{ transform: `translateX(-${100 - progress}%)` }}
      />
    </Progress.Root>
  );
}
