import Image from "next/image";
import {Box, Link, Typography} from "@mui/material";

export default function Home() {
  return (
      <Box
          sx={{
              display: 'flex',
              minHeight: '100vh',
              alignItems: 'center',
              justifyContent: 'center',
          }}
      >
          <Box
              component="main"
              sx={{
                  display: 'flex',
                  minHeight: '100vh',
                  width: '100%',
                  maxWidth: '768px',
                  flexDirection: 'column',
                  alignItems: 'center',
                  justifyContent: 'space-between',
                  py: 16,
                  px: 8,
                  '@media (min-width: 600px)': {
                      alignItems: 'flex-start',
                  },
              }}
          >
        <Image
          src="/nextjs/next.svg"
          alt="Next.js logo"
          width={100}
          height={20}
          priority
        />
              <Box
                  sx={{
                      display: 'flex',
                      flexDirection: 'column',
                      alignItems: 'center',
                      gap: 3,
                      textAlign: 'center',
                      '@media (min-width: 600px)': {
                          alignItems: 'flex-start',
                          textAlign: 'left',
                      },
                  }}
              >
                  <Typography
                      variant="h3"
                      sx={{
                          maxWidth: '20rem',
                          fontWeight: 600,
                          lineHeight: 1.25,
                      }}
                  >
            To get started, edit the page.tsx file.
                  </Typography>
                  <Typography
                      variant="body1"
                      sx={{
                          maxWidth: '28rem',
                          fontSize: '1.125rem',
                          lineHeight: 1.75,
                          color: 'text.secondary',
                      }}
                  >
            Looking for a starting point or more instructions? Head over to{" "}
                      <Link
              href="https://vercel.com/templates?framework=next.js&utm_source=create-next-app&utm_medium=appdir-template-tw&utm_campaign=create-next-app"
              sx={{fontWeight: 500}}
            >
              Templates
                      </Link>{" "}
            or the{" "}
                      <Link
              href="https://nextjs.org/learn?utm_source=create-next-app&utm_medium=appdir-template-tw&utm_campaign=create-next-app"
              sx={{fontWeight: 500}}
            >
              Learning
                      </Link>{" "}
            center.
                  </Typography>
              </Box>
              <Box
                  sx={{
                      display: 'flex',
                      flexDirection: 'column',
                      gap: 2,
                      '@media (min-width: 600px)': {
                          flexDirection: 'row',
                      },
                  }}
              >
                  <Link
                      sx={{
                          display: 'flex',
                          height: '48px',
                          width: '100%',
                          alignItems: 'center',
                          justifyContent: 'center',
                          gap: 1,
                          borderRadius: '24px',
                          bgcolor: 'primary.main',
                          px: 2.5,
                          color: 'white',
                          textDecoration: 'none',
                          '&:hover': {
                              bgcolor: 'primary.dark',
                          },
                          '@media (min-width: 960px)': {
                              width: '158px',
                          },
                      }}
            href="https://vercel.com/new?utm_source=create-next-app&utm_medium=appdir-template-tw&utm_campaign=create-next-app"
            target="_blank"
            rel="noopener noreferrer"
          >
            <Image
              src="/nextjs/vercel.svg"
              alt="Vercel logomark"
              width={16}
              height={16}
            />
            Deploy Now
                  </Link>
                  <Link
                      sx={{
                          display: 'flex',
                          height: '48px',
                          width: '100%',
                          alignItems: 'center',
                          justifyContent: 'center',
                          borderRadius: '24px',
                          border: 1,
                          borderColor: 'divider',
                          px: 2.5,
                          textDecoration: 'none',
                          color: 'text.primary',
                          '&:hover': {
                              borderColor: 'transparent',
                              bgcolor: 'action.hover',
                          },
                          '@media (min-width: 960px)': {
                              width: '158px',
                          },
                      }}
            href="https://nextjs.org/docs?utm_source=create-next-app&utm_medium=appdir-template-tw&utm_campaign=create-next-app"
            target="_blank"
            rel="noopener noreferrer"
          >
            Documentation
                  </Link>
              </Box>
          </Box>
      </Box>
  );
}
