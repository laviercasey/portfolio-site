export default function HabrIcon({ className }: { className?: string }) {
  return (
    <svg
      viewBox="0 0 24 24"
      fill="none"
      stroke="currentColor"
      strokeWidth={1.5}
      strokeLinecap="round"
      strokeLinejoin="round"
      className={className}
    >
      <path d="M2 4h4v8H2V4zm0 8h4v8H2v-8zM10 4h4v16h-4V4zM18 4h4v8h-4V4zm0 8h4v8h-4v-8z" />
    </svg>
  );
}
