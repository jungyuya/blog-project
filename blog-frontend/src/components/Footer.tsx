// blog-frontend/src/components/Footer.tsx
export default function Footer() {
  const year = new Date().getFullYear();
  return (
    <footer className="bg-gray-800 text-white p-4 mt-8 text-center">
      <div className="container mx-auto">
        <p>&copy; {year} LEE JUNGYU BLOG. All rights reserved.</p>
      </div>
    </footer>
  );
}