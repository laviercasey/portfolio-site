import dynamic from 'next/dynamic';
import { setRequestLocale } from 'next-intl/server';
import { projectsService } from '@/entities/project';
import { contentService } from '@/entities/content';
import { careerService } from '@/entities/career';
import { HeroSection } from '@/widgets/hero-section';
import { Marquee } from '@/shared/ui';

const AboutSection = dynamic(() => import('@/widgets/about-section/ui/AboutSection'));
const ProjectCarousel = dynamic(() => import('@/widgets/project-carousel/ui/ProjectCarousel'));
const HowIWorkSection = dynamic(() => import('@/widgets/how-i-work-section/ui/HowIWorkSection'));
const EducationSnapshot = dynamic(() => import('@/widgets/education-snapshot/ui/EducationSnapshot'));
const WorkHistoryTimeline = dynamic(() => import('@/widgets/work-history-timeline/ui/WorkHistoryTimeline'));
const ContactCTA = dynamic(() => import('@/widgets/contact-cta/ui/ContactCTA'));

export const revalidate = 600;

export function generateStaticParams() {
  return [];
}

export default async function HomePage({ params }: { params: Promise<{ locale: string }> }) {
  const { locale } = await params;
  setRequestLocale(locale);

  const [content, allProjects, career] = await Promise.all([
    contentService.getHomepage(),
    projectsService.list(),
    careerService.getAll(),
  ]);

  const projects = allProjects
    .filter((p) => p.featured)
    .sort((a, b) => a.order - b.order);

  const vis = content.visibility;
  const marqueeItems = content.marqueeItems?.length
    ? content.marqueeItems
    : ['Full-Stack Dev', 'Data Analysis', 'Python', 'TypeScript', 'Research', 'UI/UX', 'SQL', 'Machine Learning', 'Next.js', 'Freelance'];

  return (
    <>
      <HeroSection content={content} />

      {vis?.showMarquee !== false && (
        <div className="border-y border-border/40 py-3 overflow-hidden">
          <Marquee
            items={marqueeItems}
            className="text-sm font-mono text-foreground/40"
            speed={50}
          />
        </div>
      )}

      <AboutSection content={content} projects={allProjects} />
      <ProjectCarousel projects={projects} visibility={vis} />
      {vis?.showHowIWork !== false && <HowIWorkSection content={content} />}
      <EducationSnapshot career={career} />
      <WorkHistoryTimeline career={career} />
      {vis?.showContactCTA !== false && <ContactCTA content={content} />}
    </>
  );
}
