import { Project } from './types';
import { getAssetPath } from '@/shared/lib/utils';

export const projects: Project[] = [
  {
    id: '1',
    title: 'Med Reminder Bot',
    shortDescription: 'Это современный Telegram-бот для напоминаний о приёме лекарств, с поддержкой подписки, гибкими напоминаниями, персональными настройками и админ-панелью.',
    imageUrl: getAssetPath('/images/projects/medreminderbot.jpg'),
    technologies: ['Python', 'aiogram', 'Postgresql', 'APScheduler', 'Telegram Payments'],
    category: 'fullstack',
    demoUrl: '',
    githubUrl: 'https://github.com/laviercasey/med-reminder-bot',
    featured: false,
  },
  {
    id: '2',
    title: 'CRM Handmade',
    shortDescription: 'CRM предназначеный для хэндмайд мастеров (в стадии разработки)',
    imageUrl: getAssetPath('/images/projects/CRM.jpg'),
    technologies: ['Vue', 'TypeScript', 'FastApi'],
    category: 'fullstack',
    demoUrl: '#',
    githubUrl: 'https://github.com/laviercasey/CRM-handmade',
    featured: false,
  },
  {
    id: '3',
    title: 'Backend telegram shop',
    shortDescription: 'Высокопроизводительный API сервис для обработки данных',
    imageUrl: getAssetPath('/images/projects/backend-tg-shop.jpg'),
    technologies: ['FastApi', 'PostgreSQL', 'JWT', 'Docker'],
    category: 'backend',
    githubUrl: 'https://github.com/laviercasey/backend-telegram-shop',
    featured: false,
  },
  {
    id: '4',
    title: 'Портфолио Сайт',
    shortDescription: 'Современный сайт-портфолио с анимациями и адаптивным дизайном',
    imageUrl: getAssetPath('/images/projects/porfolio-site.jpg'),
    technologies: ['Next.js', 'React', 'TypeScript', 'Tailwind CSS', 'Framer Motion'],
    category: 'frontend',
    githubUrl: 'https://github.com/laviercasey/portfolio-site',
    featured: true,
  },
];
