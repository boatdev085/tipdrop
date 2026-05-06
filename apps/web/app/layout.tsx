import type { Metadata } from "next";
import Link from "next/link";
import "./globals.css";

export const metadata: Metadata = {
  title: "TipDrop",
  description: "Flutter-first tipping for service workers."
};

export default function RootLayout({
  children
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en">
      <body>
        <div className="shell">
          <header className="nav">
            <Link className="brand" href="/">
              TipDrop
            </Link>
            <nav className="links" aria-label="Main navigation">
              <Link href="/leaderboard">Leaderboard</Link>
              <Link href="/discover">Discover</Link>
              <Link href="/legal/privacy">Privacy</Link>
              <Link href="/legal/terms">Terms</Link>
            </nav>
          </header>
          {children}
        </div>
      </body>
    </html>
  );
}
