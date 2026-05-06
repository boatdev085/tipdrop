type PublicProfilePreviewPageProps = {
  params: {
    username: string;
  };
};

export default function PublicProfilePreviewPage({
  params
}: PublicProfilePreviewPageProps) {
  const username = decodeURIComponent(params.username);

  return (
    <main className="page">
      <p className="eyebrow">Public profile</p>
      <h1 className="title">{username.startsWith("@") ? username : "Worker preview"}</h1>
      <p className="copy">
        This route is reserved for public worker profile previews. Payment
        details must never be exposed here.
      </p>
    </main>
  );
}
