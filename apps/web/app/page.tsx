export default function HomePage() {
  return (
    <main className="page">
      <p className="eyebrow">Support web</p>
      <h1 className="title">TipDrop</h1>
      <p className="copy">
        Public support pages for leaderboard, discovery, worker profile previews,
        SEO, privacy, and terms. Core tipping stays in the Flutter app first.
      </p>
      <section className="list" aria-label="Support surfaces">
        <article className="item">
          <h2>Leaderboard</h2>
          <p>Public ranking backed by the Go API.</p>
        </article>
        <article className="item">
          <h2>Discover</h2>
          <p>Browse public worker profiles without exposing payment data.</p>
        </article>
        <article className="item">
          <h2>Profiles</h2>
          <p>Shareable previews that can deep link into the app.</p>
        </article>
      </section>
    </main>
  );
}
