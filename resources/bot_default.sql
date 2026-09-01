-- ==============================================================================
-- 1. THE NODES (The States)
-- ==============================================================================
-- Main Navigation Nodes (Centered & Top)
INSERT INTO bot_nodes (id, node_type, is_start, is_locked, pos_x, pos_y) VALUES 
("1", 'message', TRUE, TRUE, 500, 0),    -- Welcome (The Root)
("2", 'message', FALSE, FALSE, -300, 200),  -- About Menu
("3", 'message', FALSE, FALSE, 260, 200),  -- News & Activities Menu
("4", 'message', FALSE, FALSE, 670, 200),  -- Student Services Menu
("5", 'message', FALSE, FALSE, 1150, 200);  -- Feedback Menu

-- Feedback Inputs (Flowing from node 5)
INSERT INTO bot_nodes (id, node_type, is_start, is_locked, pos_x, pos_y) VALUES 
("6", 'input', FALSE, TRUE, 1000, 350),   -- Org Feedback Input
("7", 'input', FALSE, TRUE, 1300, 350),   -- Tech Feedback Input
("8", 'input', FALSE, TRUE, 1000, 650),  -- Bug Report Input
("9", 'input', FALSE, TRUE, 1300, 650),  -- General/Other Input
("10", 'message', FALSE, TRUE, 1150, 950); -- Thank You Message (The End)

-- Redirect Actions
INSERT INTO bot_nodes (id, node_type, is_start, is_locked, pos_x, pos_y) VALUES 
("11", 'action', FALSE, FALSE, -450, 400),     -- Redirect: Identity
("12", 'action', FALSE, FALSE, -450, 550),   -- Redirect: Structure
("13", 'action', FALSE, FALSE, -150, 400),   -- Redirect: Council
("14", 'action', FALSE, FALSE, -150, 550),   -- Redirect: Join Us
("15", 'action', FALSE, FALSE, 110, 400),   -- Redirect: Events
("16", 'action', FALSE, FALSE, 260, 550),   -- Redirect: News
("17", 'action', FALSE, FALSE, 410, 400),   -- Redirect: Donations
("18", 'action', FALSE, FALSE, 670, 400);   -- Redirect: Profile

-- Feedback Actions (LOCKED - Processing logic between Input and Thank You)
INSERT INTO bot_nodes (id, node_type, is_start, is_locked, pos_x, pos_y) VALUES 
("19", 'action', FALSE, TRUE, 1000, 500),   -- Action: Process Org
("20", 'action', FALSE, TRUE, 1300, 500),   -- Action: Process Tech
("21", 'action', FALSE, TRUE, 1000, 800),  -- Action: Process Bug
("22", 'action', FALSE, TRUE, 1300, 800);  -- Action: Process General


-- ==============================================================================
-- 2. NODE TRANSLATIONS (Content)
-- ==============================================================================
INSERT INTO bot_node_translations (node_id, language, content) VALUES 
-- Node 1: Welcome
("1", 'en', 'Welcome to the Engineering Association! I am your digital assistant. How can I help you today?'),
("1", 'ar', 'مرحباً بك في الجمعية الهندسية! أنا المساعد الرقمي الخاص بك. كيف يمكنني مساعدتك اليوم؟'),

-- Node 2: About Menu
("2", 'en', 'We are the legitimate body representing all students of the Faculty of Engineering. What would you like to know more about?'),
("2", 'ar', 'نحن الجسم الشرعي الممثل لجميع طلاب وطالبات كلية الهندسة في جامعة الخرطوم. ماذا تود أن تعرف عنا؟'),

-- Node 3: Activities Menu
("3", 'en', 'Stay updated with our latest workshops, student issues, and solidarity initiatives. Choose a category:'),
("3", 'ar', 'ابق على اطلاع بأحدث الورش التدريبية، قضايا الطلاب، ومبادرات التكافل. اختر فئة:'),

-- Node 4: Services
("4", 'en', 'Access your student portal to view notifications or check certificates earned from our workshops and courses.'),
("4", 'ar', 'قم بالوصول إلى بوابتك الطلابية لعرض الإشعارات أو التحقق من الشهادات المكتسبة من ورشنا ودوراتنا.'),

-- Node 5: Feedback Menu
("5", 'en', 'Your voice matters. What kind of feedback would you like to provide?'),
("5", 'ar', 'صوتك يهمنا. ما نوع الملاحظات التي تود تقديمها؟'),

-- Nodes 6-9: Inputs
("6", 'en', 'Please type your feedback regarding the Association''s organization and events:'),
("6", 'ar', 'يرجى كتابة ملاحظاتك حول تنظيم الجمعية وفعالياتها:'),
("7", 'en', 'Please describe the technical request or suggestion for the Technical Office:'),
("7", 'ar', 'يرجى وصف طلبك التقني أو مقترحك للمكتب التقني:'),
("8", 'en', 'Please describe the bug or issue you encountered on the website:'),
("8", 'ar', 'يرجى وصف العطل أو المشكلة التي واجهتها في الموقع الإلكتروني:'),
("9", 'en', 'Please type your general feedback or inquiry:'),
("9", 'ar', 'يرجى كتابة ملاحظاتك العامة أو استفسارك:'),

-- Node 10: Thank You
("10", 'en', 'Thank you! Your submission has been recorded. What would you like to do next?'),
("10", 'ar', 'شكراً لك! تم حفظ مشاركتك بنجاح. ماذا تود أن تفعل بعد ذلك؟'),

-- Node 11: Redirect to Identity
("11", 'en', 'Taking you to our core identity. Here is where our legitimacy and values as a student body live.'),
("11", 'ar', 'نأخذك الآن إلى جوهر هويتنا. هنا حيث تعيش قيمنا وشرعيتنا كجسم طلابي يمثلك.'),

-- Node 12: Redirect to Structure (The 8 Secretariats)
("12", 'en', 'Exploring the 8 secretariats—the heart and brain of our operations working for you.'),
("12", 'ar', 'جاري استعراض الأمانات الثمانية؛ القلب النابض والعقل المدبر لعملياتنا في خدمتك.'),

-- Node 13: Redirect to Council (The Treasure)
("13", 'en', 'Opening the gates to the Thirty Council: The supreme legislative and oversight authority of the students.'),
("13", 'ar', 'نفتح لك الآن أبواب المجلس الثلاثيني: السلطة التشريعية والرقابية العليا لطلاب الكلية.'),

-- Node 14: Redirect to Join Us
("14", 'en', 'Ready to be part of the change? Leading you to the application portal...'),
("14", 'ar', 'هل أنت مستعد لتكون جزءاً من التغيير؟ نوجهك الآن إلى بوابة الانضمام...'),

-- Node 15: Redirect to Events
("15", 'en', 'Check out our latest workshops and courses. Knowledge is our most valuable asset!'),
("15", 'ar', 'اكتشف آخر ورشنا ودوراتنا التدريبية. المعرفة هي أغلى ما نملك!'),

-- Node 16: Redirect to News
("16", 'en', 'Stay updated! Heading to the news and announcements section.'),
("16", 'ar', 'ابقَ على اطلاع! ننتقل الآن إلى قسم الأخبار والإعلانات الرسمية.'),

-- Node 17: Redirect to Donations
("17", 'en', 'Moving to our social solidarity initiatives. Together, we support one another.'),
("17", 'ar', 'ننتقل إلى مبادرات التكافل الاجتماعي. معاً، يدٌ واحدة تدعم الجميع.'),

-- Node 18: Redirect to Profile
("18", 'en', 'Accessing your student portal. Here you can view your notifications and earned certificates.'),
("18", 'ar', 'نفتح بوابتك الطلابية. هنا يمكنك متابعة إشعاراتك والتحقق من شهاداتك المكتسبة.'),

-- Node 19: Action Process - Org Feedback
("19", 'en', 'Sending your organizational feedback to the General Secretariat for review...'),
("19", 'ar', 'جاري إرسال ملاحظاتك التنظيمية إلى الأمانة العامة للمراجعة...'),

-- Node 20: Action Process - Technical Feedback
("20", 'en', 'Submitting your request directly to the Technical Office team. We are on it!'),
("20", 'ar', 'يتم الآن تسليم طلبك مباشرة إلى فريق المكتب التقني. نحن نعمل عليه!'),

-- Node 21: Action Process - Bug Report
("21", 'en', 'Reporting the issue to our developers. Thank you for helping us improve the digital arm!'),
("21", 'ar', 'جاري إبلاغ المطورين بالمشكلة. شكراً لمساعدتك في تحسين ذراعنا الرقمي!'),

-- Node 22: Action Process - General Feedback
("22", 'en', 'Processing your general feedback. We appreciate your transparency and commitment.'),
("22", 'ar', 'جاري معالجة ملاحظاتك العامة. نحن نقدر شفافيتك والتزامك معنا.');

-- ==============================================================================
-- 3. THE ACTIONS MAP (Ties DB to Go Logic)
-- ==============================================================================
-- Redirects
INSERT INTO bot_actions (node_id, action_type, action_text) VALUES 
("11", 'redirect', '/about/association'),
("12", 'redirect', '/about/organization-structure'), 
("13", 'redirect', '/about/council-of-thirty'),
("14", 'redirect', '/login'),
("15", 'redirect', '/events'),
("16", 'redirect', '/posts/news'),
("17", 'redirect', '/posts/donations'),
("18", 'redirect', '/profile');

-- Feedback Processing (action_text is empty because Go handles the logic)
INSERT INTO bot_actions (node_id, action_type, action_text) VALUES 
("19", 'feedback', 'organization_feedback'),
("20", 'feedback', 'technical_feedback'),
("21", 'feedback', 'bug_report'),
("22", 'feedback', 'general_feedback');


-- ==============================================================================
-- 4. THE EDGES (The Routing Logic)
-- ==============================================================================
INSERT INTO bot_edges (id, from_node_id, to_node_id, keyword) VALUES 
-- From Welcome (1)
("1", "1", "2", 'go_about'),
("2", "1", "3", 'go_activities'),
("3", "1", "4", 'go_services'),
("4", "1", "5", 'go_feedback'),

-- From About (2)
("5", "2", "11", 'read_identity'),
("6", "2", "12", 'read_structure'),
("7", "2", "13", 'read_council'),
("8", "2", "14", 'join_us'),
("9", "2", "1", 'back_main'),

-- From Activities (3)
("10", "3", "15", 'view_events'),
("11", "3", "16", 'view_news'),
("12", "3", "17", 'view_donations'),
("13", "3", "1", 'back_main'),

-- From Services (4)
("14", "4", "18", 'view_profile'),
("15", "4", "1", 'back_main'),

-- From Feedback Menu (5)
("16", "5", "6", 'org_feedback'),
("17", "5", "7", 'tech_feedback'),
("18", "5", "8", 'bug_feedback'),
("19", "5", "9", 'gen_feedback'),
("20", "5", "1", 'back_main'),

-- Process Feedback Inputs -> Actions
("21", "6", "19", 'submit_org'),
("22", "7", "20", 'submit_tech'),
("23", "8", "21", 'submit_bug'),
("24", "9", "22", 'submit_gen'),

-- From feedback to Thank you

("25", "19", "10", 'finish_org'),
("26", "20", "10", 'finish_tech'),
("27", "21", "10", 'finish_bug'),
("28", "22", "10", 'finish_gen'),

-- Return from Thank You (10)
("29", "10", "1", 'back_main');


-- ==============================================================================
-- 5. EDGE TRANSLATIONS (Button Labels)
-- ==============================================================================
INSERT INTO bot_edge_translations (edge_id, language, label) VALUES 
-- Welcome Menu Buttons
("1", 'en', 'About Us'), ("1", 'ar', 'من نحن'),
("2", 'en', 'News & Events'), ("2", 'ar', 'الأخبار والفعاليات'),
("3", 'en', 'Student Services'), ("3", 'ar', 'الخدمات الطلابية'),
("4", 'en', 'Feedback & Support'), ("4", 'ar', 'المقترحات والدعم'),

-- About Menu Buttons
("5", 'en', 'Our Identity'), ("5", 'ar', 'هويتنا'),
("6", 'en', 'Secretariats Structure'), ("6", 'ar', 'هيكل الأمانات'),
("7", 'en', 'The Thirty Council'), ("7", 'ar', 'المجلس الثلاثيني'),
("8", 'en', 'Join the Association'), ("8", 'ar', 'انضم إلينا'),
("9", 'en', 'Back to Main Menu'), ("9", 'ar', 'العودة للقائمة الرئيسية'),

-- Activities Menu Buttons
("10", 'en', 'Upcoming Events'), ("10", 'ar', 'الفعاليات القادمة'),
("11", 'en', 'Latest News'), ("11", 'ar', 'آخر الأخبار'),
("12", 'en', 'Social Solidarity'), ("12", 'ar', 'التكافل الاجتماعي'),
("13", 'en', 'Back to Main Menu'), ("13", 'ar', 'العودة للقائمة الرئيسية'),

-- Services Menu Buttons
("14", 'en', 'My Profile & Certificates'), ("14", 'ar', 'ملفي الشخصي والشهادات'),
("15", 'en', 'Back to Main Menu'), ("15", 'ar', 'العودة للقائمة الرئيسية'),

-- Feedback Menu Buttons
("16", 'en', 'Organizational Feedback'), ("16", 'ar', 'ملاحظات تنظيمية'),
("17", 'en', 'Technical Office Request'), ("17", 'ar', 'طلب للمكتب التقني'),
("18", 'en', 'Report a Website Bug'), ("18", 'ar', 'الإبلاغ عن عطل بالموقع'),
("19", 'en', 'General Feedback'), ("19", 'ar', 'ملاحظات عامة'),
("20", 'en', 'Back to Main Menu'), ("20", 'ar', 'العودة للقائمة الرئيسية'),

-- Submit buttons for inputs
("21", 'en', 'Submit'), ("21", 'ar', 'إرسال'),
("22", 'en', 'Submit'), ("22", 'ar', 'إرسال'),
("23", 'en', 'Submit'), ("23", 'ar', 'إرسال'),
("24", 'en', 'Submit'), ("24", 'ar', 'إرسال'),

-- Post-Submit to Thank you
("25", 'en', 'Continue'), ("25", 'ar', 'متابعة'),
("26", 'en', 'Continue!'), ("26", 'ar', 'متابعة'),
("27", 'en', 'Continue!'), ("27", 'ar', 'متابعة'),
("28", 'en', 'Continue!'), ("28", 'ar', 'متابعة'),

-- Post-Action return
("29", 'en', 'Return to Main Menu'), ("29", 'ar', 'العودة للقائمة الرئيسية');