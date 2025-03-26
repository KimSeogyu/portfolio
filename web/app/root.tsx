import { LinksFunction, LoaderFunction, json } from "@remix-run/node";
import {
  Links,
  Meta,
  Outlet,
  Scripts,
  ScrollRestoration,
  useLoaderData,
} from "@remix-run/react";
import tailwindStyles from "~/tailwind.css";

export const links: LinksFunction = () => [
  { rel: "stylesheet", href: tailwindStyles },
];

export const loader: LoaderFunction = async () => {
  return json({
    ENV: {
      API_URL: process.env.API_URL || "http://localhost:8080",
    },
  });
};

export default function App() {
  const data = useLoaderData<typeof loader>();

  return (
    <html lang="ko">
      <head>
        <meta charSet="utf-8" />
        <meta name="viewport" content="width=device-width, initial-scale=1" />
        <Meta />
        <Links />
      </head>
      <body className="bg-gray-100 min-h-screen">
        <header className="bg-white shadow">
          <nav className="max-w-7xl mx-auto px-4 py-4 sm:px-6 lg:px-8">
            <div className="flex justify-between">
              <a href="/" className="text-xl font-bold text-gray-800">
                포트폴리오
              </a>
              <div className="flex space-x-4">
                <a href="/board" className="text-gray-700 hover:text-blue-500">
                  게시판
                </a>
                <a href="/login" className="text-gray-700 hover:text-blue-500">
                  로그인
                </a>
              </div>
            </div>
          </nav>
        </header>
        <main className="max-w-7xl mx-auto py-6 sm:px-6 lg:px-8">
          <Outlet />
        </main>
        <footer className="bg-white border-t mt-12 py-6">
          <div className="max-w-7xl mx-auto px-4 text-center text-gray-600 text-sm">
            © {new Date().getFullYear()} 포트폴리오 프로젝트
          </div>
        </footer>
        <script
          dangerouslySetInnerHTML={{
            __html: `window.ENV = ${JSON.stringify(data.ENV)}`,
          }}
        />
        <ScrollRestoration />
        <Scripts />
      </body>
    </html>
  );
} 