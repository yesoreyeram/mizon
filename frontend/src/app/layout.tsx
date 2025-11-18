import type { Metadata } from "next";
import "./globals.css";

export const metadata: Metadata = {
  title: "Mizon - E-commerce Platform",
  description: "Minimal Amazon-like e-commerce platform",
};

export default function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <html lang="en">
      <body>{children}</body>
    </html>
  );
}
