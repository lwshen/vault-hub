'use client';

import { motion } from 'framer-motion';
import HeroGeometric from './hero-geometric';
import './hero-section.css';
export default function HeroSection({
  badge = 'VaultHub',
  title1 = 'Secure Your',
  title2 = 'Configuration',
}: {
  badge?: string;
  title1?: string;
  title2?: string;
}) {
  const fadeUpVariants = {
    hidden: { opacity: 0, y: 30 },
    visible: (i: number) => ({
      opacity: 1,
      y: 0,
      transition: {
        duration: 1,
        delay: 0.5 + i * 0.2,
        ease: [0.25, 0.4, 0.25, 1],
      },
    }),
  };

  return (
    <div className="relative min-h-screen w-full flex items-center justify-center overflow-hidden bg-background pt-16">
      <div className="absolute inset-0 bg-gradient-to-br from-emerald-500/[0.05] via-transparent to-blue-500/[0.05] blur-3xl dark:from-emerald-500/[0.05] dark:to-blue-500/[0.05]" />

      <div className="absolute inset-0 overflow-hidden">
        <HeroGeometric
          delay={0.3}
          width={600}
          height={140}
          rotate={12}
          gradient="from-emerald-500/[0.15]"
          className="left-[-10%] md:left-[-5%] top-[15%] md:top-[20%]"
        />

        <HeroGeometric
          delay={0.5}
          width={500}
          height={120}
          rotate={-15}
          gradient="from-blue-500/[0.15]"
          className="right-[-5%] md:right-[0%] top-[70%] md:top-[75%]"
        />

        <HeroGeometric
          delay={0.4}
          width={300}
          height={80}
          rotate={-8}
          gradient="from-violet-500/[0.15]"
          className="left-[5%] md:left-[10%] bottom-[5%] md:bottom-[10%]"
        />

        <HeroGeometric
          delay={0.6}
          width={200}
          height={60}
          rotate={20}
          gradient="from-amber-500/[0.15]"
          className="right-[15%] md:right-[20%] top-[10%] md:top-[15%]"
        />

        <HeroGeometric
          delay={0.7}
          width={150}
          height={40}
          rotate={-25}
          gradient="from-cyan-500/[0.15]"
          className="left-[20%] md:left-[25%] top-[5%] md:top-[10%]"
        />
      </div>

      <div className="relative z-10 container mx-auto px-4 md:px-6">
        <div className="max-w-3xl mx-auto text-center">
          <motion.div
            custom={0}
            variants={fadeUpVariants}
            initial="hidden"
            animate="visible"
            className="inline-flex items-center gap-2 px-3 py-1 rounded-full bg-muted border border-border mb-8 md:mb-12"
          >
            <div className="flex items-center justify-center w-5 h-5 bg-emerald-500/20 rounded-full">
              <svg
                xmlns="http://www.w3.org/2000/svg"
                width="12"
                height="12"
                viewBox="0 0 24 24"
                fill="none"
                stroke="currentColor"
                strokeWidth="2"
                strokeLinecap="round"
                strokeLinejoin="round"
                className="text-emerald-400"
              >
                <rect width="18" height="11" x="3" y="11" rx="2" ry="2" />
                <path d="M7 11V7a5 5 0 0 1 10 0v4" />
              </svg>
            </div>
            <span className="text-sm text-muted-foreground tracking-wide">{badge}</span>
          </motion.div>

          <motion.div custom={1} variants={fadeUpVariants} initial="hidden" animate="visible">
            <h1 className="text-4xl sm:text-6xl md:text-8xl font-bold mb-6 md:mb-8 tracking-tight">
              <span className="bg-clip-text text-transparent bg-gradient-to-b from-foreground to-foreground/80 dark:from-white dark:to-white/80">{title1}</span>
              <br />
              <span
                className="bg-clip-text text-transparent bg-gradient-to-r from-emerald-500 via-primary to-blue-500 dark:from-emerald-300 dark:via-white/90 dark:to-blue-300 pacifico-font"
              >
                {title2}
              </span>
            </h1>
          </motion.div>

          <motion.div
            custom={2}
            variants={fadeUpVariants}
            initial="hidden"
            animate="visible"
            className="mt-10 md:mt-14 max-w-2xl mx-auto"
          >
            <p className="text-base sm:text-lg md:text-xl text-foreground/70 leading-relaxed font-light">
              Keep your <span className="text-emerald-500 font-medium">API keys</span> and{' '}
              <span className="text-blue-500 font-medium">secrets</span> safe and easy to manage
            </p>
          </motion.div>
        </div>
      </div>

      <div className="absolute inset-0 bg-gradient-to-t from-background via-transparent to-background/80 pointer-events-none" />
    </div>
  );
}
