"use client";

import { motion } from 'framer-motion';
import Image from 'next/image';
import { AboutMe } from '@/widgets/AboutMe';
import { Button } from '@/shared/ui/Button';
import { profile } from '@/entities/Profile';

export default function AboutPage() {
  return (
    <div className="pt-24 relative min-h-screen bg-gradient-to-br from-slate-50 via-blue-50 to-indigo-100 dark:from-gray-900 dark:via-gray-800 dark:to-indigo-900">
      <div className="absolute inset-0 overflow-hidden pointer-events-none">
        <div className="absolute top-20 left-10 w-72 h-72 bg-blue-200 dark:bg-blue-800 rounded-full mix-blend-multiply dark:mix-blend-screen filter blur-xl opacity-30 animate-pulse"></div>
        <div className="absolute top-40 right-20 w-96 h-96 bg-purple-200 dark:bg-purple-800 rounded-full mix-blend-multiply dark:mix-blend-screen filter blur-xl opacity-20 animate-pulse" style={{animationDelay: '2s'}}></div>
      </div>
      
      <section className="py-12 relative z-10">
        <div className="container mx-auto px-4 md:px-6">
          <div className="grid grid-cols-1 lg:grid-cols-2 gap-12 items-center">
            <motion.div
              initial={{ opacity: 0, x: -30 }}
              animate={{ opacity: 1, x: 0 }}
              transition={{ duration: 0.6 }}
            >
              <h1 className="text-4xl md:text-5xl font-bold mb-6">
                Обо мне
              </h1>
              <p className="text-xl text-gray-600 dark:text-gray-300 mb-6">
                Привет! Я {profile.name}, {profile.role} с более чем 3-летним опытом создания современных веб-приложений.
              </p>
              <p className="text-gray-600 dark:text-gray-300 mb-6">
                {profile.bio.full}
              </p>
              <div className="flex flex-wrap gap-4">
                <Button href="/contact" variant="outline">Связаться со мной</Button>
              </div>
            </motion.div>
            
            <motion.div
              initial={{ opacity: 0, scale: 0.9 }}
              animate={{ opacity: 1, scale: 1 }}
              transition={{ duration: 0.6, delay: 0.2 }}
              className="relative w-full aspect-square max-w-md mx-auto"
            >
              <Image
                src="/portfolio-site/images/profile.jpg"
                alt={profile.name}
                fill
                className="object-cover rounded-2xl shadow-2xl"
                sizes="(max-width: 768px) 100vw, 500px"
              />
            </motion.div>
          </div>
        </div>
      </section>
      
      <AboutMe />
      
      <section className="py-20 bg-gray-50 dark:bg-gray-900">
        <div className="container mx-auto px-4 md:px-6">
          <motion.div
            className="max-w-4xl mx-auto text-center"
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ duration: 0.6 }}
          >
            <h2 className="text-3xl md:text-4xl font-bold mb-6">Мой подход к работе</h2>
            <p className="text-xl text-gray-600 dark:text-gray-300 mb-8">
              Я верю, что лучшие проекты создаются на основе тесного сотрудничества, четкого общения и внимания к деталям.
            </p>
            
            <div className="grid grid-cols-1 md:grid-cols-3 gap-8">
              <motion.div
                className="bg-white dark:bg-gray-800 p-6 rounded-xl shadow-lg"
                whileHover={{ y: -10, transition: { duration: 0.2 } }}
              >
                <div className="text-blue-600 dark:text-blue-400 mb-4">
                  <svg
                    xmlns="http://www.w3.org/2000/svg"
                    className="h-12 w-12 mx-auto"
                    fill="none"
                    viewBox="0 0 24 24"
                    stroke="currentColor"
                  >
                    <path
                      strokeLinecap="round"
                      strokeLinejoin="round"
                      strokeWidth={2}
                      d="M9.663 17h4.673M12 3v1m6.364 1.636l-.707.707M21 12h-1M4 12H3m3.343-5.657l-.707-.707m2.828 9.9a5 5 0 117.072 0l-.548.547A3.374 3.374 0 0014 18.469V19a2 2 0 11-4 0v-.531c0-.895-.356-1.754-.988-2.386l-.548-.547z"
                    />
                  </svg>
                </div>
                <h3 className="text-xl font-bold mb-2">Инновации</h3>
                <p className="text-gray-600 dark:text-gray-300">
                  Я всегда ищу инновационные решения сложных проблем, используя современные технологии и подходы.
                </p>
              </motion.div>
              
              <motion.div
                className="bg-white dark:bg-gray-800 p-6 rounded-xl shadow-lg"
                whileHover={{ y: -10, transition: { duration: 0.2 } }}
              >
                <div className="text-blue-600 dark:text-blue-400 mb-4">
                  <svg
                    xmlns="http://www.w3.org/2000/svg"
                    className="h-12 w-12 mx-auto"
                    fill="none"
                    viewBox="0 0 24 24"
                    stroke="currentColor"
                  >
                    <path
                      strokeLinecap="round"
                      strokeLinejoin="round"
                      strokeWidth={2}
                      d="M13 10V3L4 14h7v7l9-11h-7z"
                    />
                  </svg>
                </div>
                <h3 className="text-xl font-bold mb-2">Производительность</h3>
                <p className="text-gray-600 dark:text-gray-300">
                  Я уделяю особое внимание производительности, создавая быстрые и отзывчивые приложения.
                </p>
              </motion.div>
              
              <motion.div
                className="bg-white dark:bg-gray-800 p-6 rounded-xl shadow-lg"
                whileHover={{ y: -10, transition: { duration: 0.2 } }}
              >
                <div className="text-blue-600 dark:text-blue-400 mb-4">
                  <svg
                    xmlns="http://www.w3.org/2000/svg"
                    className="h-12 w-12 mx-auto"
                    fill="none"
                    viewBox="0 0 24 24"
                    stroke="currentColor"
                  >
                    <path
                      strokeLinecap="round"
                      strokeLinejoin="round"
                      strokeWidth={2}
                      d="M12 6V4m0 2a2 2 0 100 4m0-4a2 2 0 110 4m-6 8a2 2 0 100-4m0 4a2 2 0 110-4m0 4v2m0-6V4m6 6v10m6-2a2 2 0 100-4m0 4a2 2 0 110-4m0 4v2m0-6V4"
                    />
                  </svg>
                </div>
                <h3 className="text-xl font-bold mb-2">Адаптивность</h3>
                <p className="text-gray-600 dark:text-gray-300">
                  Мои решения адаптивны и гибки, обеспечивая отличный пользовательский опыт на любых устройствах.
                </p>
              </motion.div>
            </div>
          </motion.div>
        </div>
      </section>
    </div>
  );
}
