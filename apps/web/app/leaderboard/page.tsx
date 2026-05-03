const workers = [
  { rank: 1, name: "Maya", score: "128 confirmed tips" },
  { rank: 2, name: "Anan", score: "104 confirmed tips" },
  { rank: 3, name: "Nok", score: "91 confirmed tips" }
];

export default function LeaderboardPage() {
  return (
    <main className="page">
      <p className="eyebrow">Leaderboard</p>
      <h1 className="title">Top workers</h1>
      <p className="copy">Stub leaderboard page for the public support web app.</p>
      <section className="list" aria-label="Leaderboard entries">
        {workers.map((worker) => (
          <article className="item" key={worker.rank}>
            <h2>
              #{worker.rank} {worker.name}
            </h2>
            <p>{worker.score}</p>
          </article>
        ))}
      </section>
    </main>
  );
}
