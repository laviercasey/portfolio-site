"use client";

import { Hero } from '@/widgets/Hero';
import { AboutMe } from '@/widgets/AboutMe';
import { ProjectList } from '@/widgets/ProjectList';
import { ContactForm } from '@/features/ContactForm';

export default function Home() {
  return (
    <div className="pt-16">
      <Hero />
      <AboutMe />
      <ProjectList />
      <section className="py-20 bg-gray-50 dark:bg-gray-900">
        <div className="container mx-auto px-4 md:px-6">
          <div className="text-center mb-12">
            <h2 className="text-3xl md:text-4xl font-bold mb-4">Связаться со мной</h2>
            <p className="text-xl text-gray-600 dark:text-gray-300 max-w-3xl mx-auto">
              У вас есть проект или идея? Давайте обсудим, как я могу помочь воплотить её в жизнь.
            </p>
          </div>
          <ContactForm />
        </div>
      </section>
    </div>
  );
}
