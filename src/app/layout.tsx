import './globals.css';
import { Header } from '@/widgets/Header';
import { Footer } from '@/widgets/Footer';
import type { Metadata } from 'next';

export const metadata: Metadata = {
  title: 'Лавьер Кейси | Портфолио',
  description: 'Fullstack разработчик с экспертизой в React, Next.js и TypeScript',
};

export default function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <html lang="ru">
      <body className="bg-gray-50 dark:bg-gray-900 text-gray-900 dark:text-gray-100">
        <div className="min-h-screen flex flex-col">
          <Header />
          <main className="flex-grow">{children}</main>
          <Footer />
        </div>
      </body>
    </html>
  );
}
