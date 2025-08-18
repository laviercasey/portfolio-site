export interface ProfileData {
  name: string;
  role: string;
  avatar: string;
  email: string;
  phone: string;
  telegram: string;
  bio: {
    short: string;
    full: string;
  };
  skills: Array<{
    name: string;
    level: number;
  }>;
  experience: Array<{
    company: string;
    position: string;
    period: string;
    description: string;
  }>;
  education: Array<{
    institution: string;
    degree: string;
    period: string;
  }>;
  social: {
    github: string;
    linkedin: string;
    telegram: string;
  };
}

import { getAssetPath } from '@/shared/lib/utils';

export const profile: ProfileData = {
  name: "Лавьер Кейси",
  role: "Fullstack Developer",
  avatar: getAssetPath("/images/profile.jpg"),
  email: "laviercasey@gmail.com",
  phone: "+7 (903) 548-7402",
  telegram: "@CaseyLav",
  bio: {
    short: "Я специализируюсь на создании современных веб-приложений с использованием React, Vue, TypeScript и FastApi. Мой фокус — на производительности, масштабируемости и удобстве использования.",
    full: "Более 3 лет я занимаюсь разработкой веб-приложений, начиная от простых лендингов и заканчивая сложными корпоративными системами. Я стремлюсь к постоянному совершенствованию своих навыков и изучению новых технологий. Мой подход к разработке основан на создании чистого, поддерживаемого кода и использовании лучших практик отрасли. Я верю, что хороший код должен быть не только функциональным, но и элегантным."
  },
  skills: [
    { name: 'Vue.js/React', level: 80 },
    { name: 'JavaScript/TypeScript', level: 80 },
    { name: 'Next.js/Nuxt.js', level: 70 },
    { name: 'FastApi', level: 70 },
    { name: 'Tailwind CSS', level: 70 },
    { name: 'Redux/Pinia/Vuex', level: 70 },
    { name: 'Redis', level: 60 },
    { name: 'PostgreSQL', level: 70 },
    { name: 'Docker', level: 65 },
  ],
  experience: [
    {
      company: 'Российская государственная детская библиотека',
      position: 'Старший группы разработки',
      period: '2023 - Present',
      description: 'Разработка и поддержка высоконагруженных веб-приложений с использованием React, Vue.js и FastAPi. Создание микросервисной/монолитной архитектуры и оптимизация производительности.'
    },
    {
      company: 'Российская государственная детская библиотека',
      position: 'Младший разработчик',
      period: '2021 - 2023',
      description: 'Разработка адаптивных веб-сайтов и приложений с использованием современных технологий. Участие в полном цикле разработки от проектирования до внедрения.'
    }
  ],
  education: [
    {
      institution: 'Московский университет имени С.Ю. Витте, Москва',
      degree: 'Прикладная информатика',
      period: '2022 - 2026'
    },
    {
      institution: 'МГОК',
      degree: 'Сетевое и системное администрирование',
      period: '2017 - 2021'
    }
  ],
  social: {
    github: 'https://github.com/laviercasey',
    linkedin: 'https://linkedin.com/in/laviercasey',
    telegram: 'https://t.me/CaseyLav'
  }
};
